package dependency

import (
	"testing"

	"github.com/TamalSaha/go-oneliners"
	"kube.ci/kubeci/apis/kubeci/v1alpha1"
)

var cleanupStep = v1alpha1.Step{
	Name:     "cleanup-step",
	Image:    "alpine",
	Commands: []string{"rm"},
	Args:     []string{"-rf", "/kubeci/*"},
}

var stepsData = [][]v1alpha1.Step{
	{
		{
			Name: "A",
		},
		{
			Name:       "B",
			Dependency: []string{"A"},
		},
		{
			Name:       "C",
			Dependency: []string{"A"},
		},
		{
			Name:       "D",
			Dependency: []string{"A", "B"},
		},
		{
			Name:       "E",
			Dependency: []string{"B", "C"},
		},
		{
			Name:       "F",
			Dependency: []string{"A"},
		},
	},
	{
		{
			Name: "A",
		},
		{
			Name:       "B",
			Dependency: []string{"A"},
		},
		{
			Name:       "C",
			Dependency: []string{"B"},
		},
		{
			Name:       "D",
			Dependency: []string{"A", "B"},
		},
		{
			Name:       "E",
			Dependency: []string{"B", "C"},
		},
		{
			Name:       "F",
			Dependency: []string{"D"},
		},
	},
}

func TestResolveDependency(t *testing.T) {
	for _, steps := range stepsData {
		if tasks, err := ResolveDependency(steps, cleanupStep); err != nil {
			t.Errorf(err.Error())
		} else {
			oneliners.PrettyJson(tasks)
		}
	}
}

func TestDagToLayers(t *testing.T) {
	for _, steps := range stepsData {
		stepsMap := make(map[string]v1alpha1.Step, 0)
		for _, step := range steps {
			stepsMap[step.Name] = step
		}

		if layers, err := dagToLayers(stepsMap); err != nil {
			t.Errorf(err.Error())
		} else {
			oneliners.PrettyJson(layers)
		}
	}
}