package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var (
	EnableStatusSubresource bool
)

func (c Workflow) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceWorkflows,
		Kind:          ResourceKindWorkflow,
		ShortNames:    []string{"wf"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "kube.ci/kubeci/apis/kubeci/v1alpha1.Workflow",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Workplan) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceWorkplans,
		Kind:          ResourceKindWorkplan,
		ShortNames:    []string{"wp"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "kube.ci/kubeci/apis/kubeci/v1alpha1.Workplan",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}