package v1alpha1

import (
	v1 "k8s.io/api/apps/v1"
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
	DeploymentName string        `json:"deploymentName"`
	Deployment     v1.Deployment `json:"deployment"`
}

// DepSvcResourceStatus is the status for a DepSvcResource
type DepSvcResourceStatus struct {
	AvailableReplicas int32               `json:"availableReplicas"`
	DeploymentStatus  v1.DeploymentStatus `json:"deploymentStatus"`
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
