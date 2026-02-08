/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/moondev/cluster-api-provider-maas/test/helpers"
)

var (
	testEnv *helpers.TestEnvironment
	ctx     = ctrl.SetupSignalHandler()
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{})
}

func TestMain(m *testing.M) {
	if !setup() {
		os.Exit(0)
	}
	defer teardown()
	code := m.Run()
	os.Exit(code)
}

// setup returns false if envtest is unavailable (e.g. missing etcd/kube-apiserver), so tests are skipped.
func setup() bool {
	utilruntime.Must(AddToScheme(scheme.Scheme))
	testEnvConfig := helpers.NewTestEnvironmentConfiguration([]string{
		path.Join("config", "crd", "bases"),
	},
	).WithWebhookConfiguration("unmanaged", path.Join("config", "webhook", "manifests.yaml"))
	var err error
	testEnv, err = testEnvConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Skipping api/v1beta1 integration tests (envtest binaries not available): %v\n", err)
		return false
	}
	if err := (&MaasCluster{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup MAASCluster webhook: %v", err))
	}
	if err := (&MaasMachine{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup MAASMachine webhook: %v", err))
	}
	if err := (&MaasMachineTemplate{}).SetupWebhookWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Unable to setup MAASMachineTemplate webhook: %v", err))
	}

	go func() {
		fmt.Println("Starting the manager")
		if err := testEnv.StartManager(ctx); err != nil {
			panic(fmt.Sprintf("Failed to start the envtest manager: %v", err))
		}
	}()
	testEnv.WaitForWebhooks()
	return true
}

func teardown() {
	if testEnv == nil {
		return
	}
	if err := testEnv.Stop(); err != nil {
		panic(fmt.Sprintf("Failed to stop envtest: %v", err))
	}
}
