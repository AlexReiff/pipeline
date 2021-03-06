// +build e2e

/*
Copyright 2019 The Tekton Authors

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

package test

import (
	"testing"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativetest "knative.dev/pkg/test"
)

const epTaskRunName = "ep-task-run"

// TestEntrypointRunningStepsInOrder is an integration test that will
// verify attempt to the get the entrypoint of a container image
// that doesn't have a cmd defined. In addition to making sure the steps
// are executed in the order specified
func TestEntrypointRunningStepsInOrder(t *testing.T) {
	c, namespace := setup(t)
	t.Parallel()

	knativetest.CleanupOnInterrupt(func() { tearDown(t, c, namespace) }, t.Logf)
	defer tearDown(t, c, namespace)

	t.Logf("Creating TaskRun in namespace %s", namespace)
	if _, err := c.TaskRunClient.Create(&v1alpha1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{Name: epTaskRunName, Namespace: namespace},
		Spec: v1alpha1.TaskRunSpec{
			TaskSpec: &v1alpha1.TaskSpec{TaskSpec: v1beta1.TaskSpec{
				Steps: []v1beta1.Step{{
					Container: corev1.Container{Image: "busybox"},
					Script:    "sleep 3 && touch foo",
				}, {
					Container: corev1.Container{Image: "ubuntu"},
					Script:    "ls foo",
				}},
			}},
		},
	}); err != nil {
		t.Fatalf("Failed to create TaskRun: %s", err)
	}

	t.Logf("Waiting for TaskRun in namespace %s to finish successfully", namespace)
	if err := WaitForTaskRunState(c, epTaskRunName, TaskRunSucceed(epTaskRunName), "TaskRunSuccess"); err != nil {
		t.Errorf("Error waiting for TaskRun to finish successfully: %s", err)
	}

}
