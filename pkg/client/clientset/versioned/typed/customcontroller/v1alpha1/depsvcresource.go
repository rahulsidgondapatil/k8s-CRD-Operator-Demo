/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/rahulsidgondapatil/sample-customController/pkg/apis/customcontroller/v1alpha1"
	scheme "github.com/rahulsidgondapatil/sample-customController/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DepSvcResourcesGetter has a method to return a DepSvcResourceInterface.
// A group's client should implement this interface.
type DepSvcResourcesGetter interface {
	DepSvcResources(namespace string) DepSvcResourceInterface
}

// DepSvcResourceInterface has methods to work with DepSvcResource resources.
type DepSvcResourceInterface interface {
	Create(*v1alpha1.DepSvcResource) (*v1alpha1.DepSvcResource, error)
	Update(*v1alpha1.DepSvcResource) (*v1alpha1.DepSvcResource, error)
	UpdateStatus(*v1alpha1.DepSvcResource) (*v1alpha1.DepSvcResource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.DepSvcResource, error)
	List(opts v1.ListOptions) (*v1alpha1.DepSvcResourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DepSvcResource, err error)
	DepSvcResourceExpansion
}

// depSvcResources implements DepSvcResourceInterface
type depSvcResources struct {
	client rest.Interface
	ns     string
}

// newDepSvcResources returns a DepSvcResources
func newDepSvcResources(c *CustomcontrollerV1alpha1Client, namespace string) *depSvcResources {
	return &depSvcResources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the depSvcResource, and returns the corresponding depSvcResource object, and an error if there is any.
func (c *depSvcResources) Get(name string, options v1.GetOptions) (result *v1alpha1.DepSvcResource, err error) {
	result = &v1alpha1.DepSvcResource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("depsvcresources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DepSvcResources that match those selectors.
func (c *depSvcResources) List(opts v1.ListOptions) (result *v1alpha1.DepSvcResourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.DepSvcResourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("depsvcresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested depSvcResources.
func (c *depSvcResources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("depsvcresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a depSvcResource and creates it.  Returns the server's representation of the depSvcResource, and an error, if there is any.
func (c *depSvcResources) Create(depSvcResource *v1alpha1.DepSvcResource) (result *v1alpha1.DepSvcResource, err error) {
	result = &v1alpha1.DepSvcResource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("depsvcresources").
		Body(depSvcResource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a depSvcResource and updates it. Returns the server's representation of the depSvcResource, and an error, if there is any.
func (c *depSvcResources) Update(depSvcResource *v1alpha1.DepSvcResource) (result *v1alpha1.DepSvcResource, err error) {
	result = &v1alpha1.DepSvcResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("depsvcresources").
		Name(depSvcResource.Name).
		Body(depSvcResource).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *depSvcResources) UpdateStatus(depSvcResource *v1alpha1.DepSvcResource) (result *v1alpha1.DepSvcResource, err error) {
	result = &v1alpha1.DepSvcResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("depsvcresources").
		Name(depSvcResource.Name).
		SubResource("status").
		Body(depSvcResource).
		Do().
		Into(result)
	return
}

// Delete takes name of the depSvcResource and deletes it. Returns an error if one occurs.
func (c *depSvcResources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("depsvcresources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *depSvcResources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("depsvcresources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched depSvcResource.
func (c *depSvcResources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DepSvcResource, err error) {
	result = &v1alpha1.DepSvcResource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("depsvcresources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}