package main

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/rahulsidgondapatil/sample-customController/conf"
	customcontroller "github.com/rahulsidgondapatil/sample-customController/pkg/apis/customcontroller/v1alpha1"
	clientset "github.com/rahulsidgondapatil/sample-customController/pkg/client/clientset/versioned"
	samplescheme "github.com/rahulsidgondapatil/sample-customController/pkg/client/clientset/versioned/scheme"
	informers "github.com/rahulsidgondapatil/sample-customController/pkg/client/informers/externalversions/customcontroller/v1alpha1"
	listers "github.com/rahulsidgondapatil/sample-customController/pkg/client/listers/customcontroller/v1alpha1"
)

const controllerAgentName = "sample-customController"

const (
	// SuccessSynced is used as part of the Event 'reason' when a DepSvcResource is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a DepSvcResource fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by DepSvcResource"
	// MessageResourceSynced is the message used for an Event fired when a DepSvcResource
	// is synced successfully
	MessageResourceSynced = "DepSvcResource synced successfully"
)

//Controller...
type Controller struct {
	Cfg *conf.Config
	// kubeclientset is a standard kubernetes clientset
	Kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	Sampleclientset clientset.Interface

	DeploymentsLister     appslisters.DeploymentLister
	DeploymentsSynced     cache.InformerSynced
	DepSvcResourceLister  listers.DepSvcResourceLister
	DepSvcResourcesSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	cfg *conf.Config,
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	depSvcResourceInformer informers.DepSvcResourceInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	runtime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		Cfg:                   cfg,
		Kubeclientset:         kubeclientset,
		Sampleclientset:       sampleclientset,
		DeploymentsLister:     deploymentInformer.Lister(),
		DeploymentsSynced:     deploymentInformer.Informer().HasSynced,
		DepSvcResourceLister:  depSvcResourceInformer.Lister(),
		DepSvcResourcesSynced: depSvcResourceInformer.Informer().HasSynced,
		workqueue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DepSvcResources"),
		recorder:              recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when DepSvcResources change
	depSvcResourceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueDepSvcResource,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueDepSvcResource(new)
		},
		DeleteFunc: controller.enqueueDepSvcResource,
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a DepSvcResource will enqueue that DepSvcResource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting DepSvcResource controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.DeploymentsSynced, c.DepSvcResourcesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process DepSvcResource resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	klog.Info("Workqueue:%+v", c.workqueue)
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// DepSvcResource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the DepSvcResource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the DepSvcResource with this namespace/name
	depSvcResource, err := c.DepSvcResourceLister.DepSvcResources(namespace).Get(name)
	if err != nil {
		// The DepSvcResource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("depSvcResource '%s' in work queue no longer exists", key))
			// We delete the corresponding service as well
			c.deleteSvc(name, namespace)
			return nil
		}

		return err
	}

	deploymentName := depSvcResource.Name
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in DepSvcResource.spec
	deployment, err := c.DeploymentsLister.Deployments(depSvcResource.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		c.newDeployNService(depSvcResource)
	}
	//klog.Info("\nsyncHandler: Created deployment:%+v\n", deployment)

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this DepSvcResource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, depSvcResource) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(depSvcResource, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the DepSvcResource resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if &depSvcResource.Spec.Replicas != nil && *depSvcResource.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("depSvcResource %s replicas: %d, deployment replicas: %d", name, depSvcResource.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.Kubeclientset.AppsV1().Deployments(depSvcResource.Namespace).Update(newDeployment(depSvcResource))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the DepSvcResource resource to reflect the
	// current state of the world
	err = c.updateDepSvcStatus(depSvcResource, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(depSvcResource, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateDepSvcStatus(depSvcResource *customcontroller.DepSvcResource, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	depSvcResourceCopy := depSvcResource.DeepCopy()
	depSvcResourceCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the DepSvcResource resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.Sampleclientset.CustomcontrollerV1alpha1().DepSvcResources(depSvcResource.Namespace).Update(depSvcResourceCopy)
	return err
}

// enqueueDepSvcResource takes a DepSvcResource resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than DepSvcResource.
func (c *Controller) enqueueDepSvcResource(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the DepSvcResource resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that DepSvcResource resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a DepSvcResource, we should not do anything more
		// with it.
		if ownerRef.Kind != "DepSvcResource" {
			return
		}

		depSvcResource, err := c.DepSvcResourceLister.DepSvcResources(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of depSvcResource '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueDepSvcResource(depSvcResource)
		return
	}
}

func (c *Controller) newDeployNService(depSvcResource *customcontroller.DepSvcResource) {
	deployment, err := c.Kubeclientset.AppsV1().Deployments(depSvcResource.Namespace).Create(newDeployment(depSvcResource))
	klog.Info("\nCreated deployment:'%v\n", deployment.Name)
	svclient := c.Kubeclientset.CoreV1().Services(depSvcResource.Namespace)
	res, err := svclient.Create(newSvcObject(deployment, c.Cfg))
	if err != nil {
		fmt.Printf("\nError:%+v", err.Error())
	}
	klog.Info("\nCreated service %v for deployment:%v\n", res.Name, deployment.Name)
}

func (c *Controller) deleteSvc(name, namespace string) {
	svclient := c.Kubeclientset.CoreV1().Services(namespace)
	svc, _ := svclient.Get(name, metav1.GetOptions{})
	if svc.ObjectMeta.Name == name {
		klog.Info("Deleting service for depSvcResource:%v...\n", name)
		deletePolicy := metav1.DeletePropagationForeground
		if err := svclient.Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		klog.Info("Deleted service for depSvcResource:%v\n", name)
	}
}

func newDeployment(depSvcResource *customcontroller.DepSvcResource) *appsv1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      depSvcResource.ObjectMeta.Name,
			Namespace: depSvcResource.ObjectMeta.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(depSvcResource, customcontroller.SchemeGroupVersion.WithKind("DepSvcResource")),
			},
		},
		Spec: v1.DeploymentSpec{
			Replicas:                depSvcResource.Spec.Replicas,
			Selector:                depSvcResource.Spec.Selector,
			Template:                depSvcResource.Spec.Template,
			Strategy:                depSvcResource.Spec.Strategy,
			MinReadySeconds:         depSvcResource.Spec.MinReadySeconds,
			RevisionHistoryLimit:    depSvcResource.Spec.RevisionHistoryLimit,
			Paused:                  depSvcResource.Spec.Paused,
			ProgressDeadlineSeconds: depSvcResource.Spec.ProgressDeadlineSeconds,
		},
	}
}

func newSvcObject(d *v1.Deployment, cfg *conf.Config) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.ObjectMeta.Name,
		},
		Spec: corev1.ServiceSpec{
			Type:     cfg.Settings.AutocreateSvcType,
			Selector: d.Spec.Selector.MatchLabels,
			Ports:    getSvcPorts(d.Spec.Template.Spec.Containers, 9090),
		},
	}
}

func getSvcPorts(containers []corev1.Container, firstPort int32) []corev1.ServicePort {
	svcPorts := make([]corev1.ServicePort, 0)
	portNumber := firstPort
	for _, c := range containers {
		svcPort := corev1.ServicePort{
			Name: c.Name,
			Port: portNumber,
			TargetPort: intstr.IntOrString{
				IntVal: c.Ports[0].ContainerPort,
			},
		}
		svcPorts = append(svcPorts, svcPort)
		portNumber++
	}
	return svcPorts
}
