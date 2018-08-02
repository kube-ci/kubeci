package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindWorkplan = "Workplan"
	ResourceWorkplans    = "workplans"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Workplan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkplanSpec   `json:"spec,omitempty"`
	Status WorkplanStatus `json:"status,omitempty"`
}

type Task struct { // analogous to a single pod
	SerialSteps   []Step // analogous to init-containers
	ParallelSteps []Step // analogous to sidecar-containers
}

type WorkplanSpec struct {
	Tasks []Task `json:"tasks,omitempty"`
	// set container environment variables from configmaps and secrets
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
}

type WorkplanStatus struct {
	Phase     string `json:"phase"`
	Reason    string `json:"reason"`
	TaskIndex int    `json:"taskIndex"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WorkplanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Workplan `json:"items"`
}