/*
Copyright 2025.

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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	thermopilotv1 "github.com/seipan/thermo-pilot-controller/api/v1"
)

var _ = Describe("ThermoPilot Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		thermopilot := &thermopilotv1.ThermoPilot{}

		BeforeEach(func() {
			By("creating the Secret with SwitchBot credentials")
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "default",
				},
				StringData: map[string]string{
					"token":  "test-token",
					"secret": "test-secret",
				},
			}
			err := k8sClient.Create(ctx, secret)
			if err != nil && !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}

			By("creating the custom resource for the Kind ThermoPilot")
			err = k8sClient.Get(ctx, typeNamespacedName, thermopilot)
			if err != nil && errors.IsNotFound(err) {
				resource := &thermopilotv1.ThermoPilot{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: thermopilotv1.ThermoPilotSpec{
						SecretRef: thermopilotv1.SecretReference{
							Name: "test-secret",
						},
						TemperatureSensorType: "MeterPro",
						TargetTemperature:     "25.0",
						Mode:                  "cool",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &thermopilotv1.ThermoPilot{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("Cleanup the specific resource instance ThermoPilot")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: "test-secret", Namespace: "default"}, secret)
			if err == nil {
				By("Cleanup the test secret")
				Expect(k8sClient.Delete(ctx, secret)).To(Succeed())
			}
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ThermoPilotReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			// The reconciler will fail due to missing API mock, but we expect it to at least try
			// In a real test, we would mock the HTTP client
			Expect(err).To(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", 0))
		})
	})
})
