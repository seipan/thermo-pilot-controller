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

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	thermopilotv1 "github.com/seipan/thermo-pilot-controller/api/v1"
)

// nolint:unused
// log is for logging in this package.
var thermopilotlog = logf.Log.WithName("thermopilot-resource")

// SetupThermoPilotWebhookWithManager registers the webhook for ThermoPilot in the manager.
func SetupThermoPilotWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&thermopilotv1.ThermoPilot{}).
		WithValidator(&ThermoPilotCustomValidator{}).
		WithDefaulter(&ThermoPilotCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-thermo-pilot-yadon3141-com-v1-thermopilot,mutating=true,failurePolicy=fail,sideEffects=None,groups=thermo-pilot.yadon3141.com,resources=thermopilots,verbs=create;update,versions=v1,name=mthermopilot-v1.kb.io,admissionReviewVersions=v1

// ThermoPilotCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ThermoPilot when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ThermoPilotCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ThermoPilotCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ThermoPilot.
func (d *ThermoPilotCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	thermopilot, ok := obj.(*thermopilotv1.ThermoPilot)

	if !ok {
		return fmt.Errorf("expected an ThermoPilot object but got %T", obj)
	}
	thermopilotlog.Info("Defaulting for ThermoPilot", "name", thermopilot.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-thermo-pilot-yadon3141-com-v1-thermopilot,mutating=false,failurePolicy=fail,sideEffects=None,groups=thermo-pilot.yadon3141.com,resources=thermopilots,verbs=create;update,versions=v1,name=vthermopilot-v1.kb.io,admissionReviewVersions=v1

// ThermoPilotCustomValidator struct is responsible for validating the ThermoPilot resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ThermoPilotCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ThermoPilotCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ThermoPilot.
func (v *ThermoPilotCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	thermopilot, ok := obj.(*thermopilotv1.ThermoPilot)
	if !ok {
		return nil, fmt.Errorf("expected a ThermoPilot object but got %T", obj)
	}
	thermopilotlog.Info("Validation for ThermoPilot upon creation", "name", thermopilot.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ThermoPilot.
func (v *ThermoPilotCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	thermopilot, ok := newObj.(*thermopilotv1.ThermoPilot)
	if !ok {
		return nil, fmt.Errorf("expected a ThermoPilot object for the newObj but got %T", newObj)
	}
	thermopilotlog.Info("Validation for ThermoPilot upon update", "name", thermopilot.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ThermoPilot.
func (v *ThermoPilotCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	thermopilot, ok := obj.(*thermopilotv1.ThermoPilot)
	if !ok {
		return nil, fmt.Errorf("expected a ThermoPilot object but got %T", obj)
	}
	thermopilotlog.Info("Validation for ThermoPilot upon deletion", "name", thermopilot.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
