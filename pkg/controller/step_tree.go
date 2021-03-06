package controller

import (
	api "github.com/kube-ci/engine/apis/engine/v1alpha1"
	"github.com/kube-ci/engine/pkg/dependency"
	v1 "k8s.io/api/core/v1"
)

func InitWorkplanTree(tasks []api.Task) [][]api.StepEntry {
	var stepTree [][]api.StepEntry
	layers := dependency.TasksToLayers(tasks)
	for _, layer := range layers {
		var stepEntries []api.StepEntry
		for _, step := range layer {
			stepEntries = append(stepEntries, api.StepEntry{
				Name:   step.Name,
				Status: api.ContainerUninitialized,
			})
		}
		stepTree = append(stepTree, stepEntries)
	}
	return stepTree
}

func UpdateWorkplanTreeForPod(stepTree [][]api.StepEntry, pod *v1.Pod) [][]api.StepEntry {
	statuses := append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...)
	containerStatuses := make(map[string]api.StepEntry)

	for _, status := range statuses {
		// status.Name = container-name
		containerStatuses[status.Name] = api.StepEntry{
			Name:           status.Name,
			PodName:        pod.Name,
			Status:         getContainerState(status.State),
			ContainerState: status.State,
		}
	}

	for i, layer := range stepTree {
		for j, step := range layer {
			stepEntry, ok := containerStatuses[step.Name]
			if ok { // only update matching steps
				stepTree[i][j] = stepEntry
			}
		}
	}

	return stepTree
}

func getContainerState(state v1.ContainerState) api.ContainerStatus {
	switch {
	case state.Running != nil:
		return api.ContainerRunning
	case state.Waiting != nil:
		return api.ContainerWaiting
	case state.Terminated != nil:
		return api.ContainerTerminated
	default:
		return api.ContainerUninitialized
	}
}
