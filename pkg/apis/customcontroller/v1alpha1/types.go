package v1alpha1

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DepSvcResource is a top-level type
type DepSvcResource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +optional
	Status DepSvcResourceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	// This is where you can define
	// your own custom spec
	Spec DepSvcResourceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// DepSvcResourceSpec is the spec for a Foo resource
type DepSvcResourceSpec struct {
	Replicas                *int32                 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	Selector                *metav1.LabelSelector  `json:"selector" protobuf:"bytes,2,opt,name=selector"`
	Template                corev1.PodTemplateSpec `json:"template" protobuf:"bytes,3,opt,name=template"`
	Strategy                v1.DeploymentStrategy  `json:"strategy,omitempty" patchStrategy:"retainKeys" protobuf:"bytes,4,opt,name=strategy"`
	MinReadySeconds         int32                  `json:"minReadySeconds,omitempty" protobuf:"varint,5,opt,name=minReadySeconds"`
	RevisionHistoryLimit    *int32                 `json:"revisionHistoryLimit,omitempty" protobuf:"varint,6,opt,name=revisionHistoryLimit"`
	Paused                  bool                   `json:"paused,omitempty" protobuf:"varint,7,opt,name=paused"`
	ProgressDeadlineSeconds *int32                 `json:"progressDeadlineSeconds,omitempty" protobuf:"varint,9,opt,name=progressDeadlineSeconds"`
}

// DepSvcResourceStatus is the status for a DepSvcResource
type DepSvcResourceStatus struct {
	ObservedGeneration  int64                    `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`
	Replicas            int32                    `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`
	UpdatedReplicas     int32                    `json:"updatedReplicas,omitempty" protobuf:"varint,3,opt,name=updatedReplicas"`
	ReadyReplicas       int32                    `json:"readyReplicas,omitempty" protobuf:"varint,7,opt,name=readyReplicas"`
	AvailableReplicas   int32                    `json:"availableReplicas,omitempty" protobuf:"varint,4,opt,name=availableReplicas"`
	UnavailableReplicas int32                    `json:"unavailableReplicas,omitempty" protobuf:"varint,5,opt,name=unavailableReplicas"`
	Conditions          []v1.DeploymentCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
	CollisionCount      *int32                   `json:"collisionCount,omitempty" protobuf:"varint,8,opt,name=collisionCount"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// no client needed for list as it's been created in above
// DepSvcResourceList is a list of DepSvcResource resources
type DepSvcResourceList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `son:"metadata,omitempty"`

	Items []DepSvcResource `json:"items"`
}
