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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ThermoPilotSpec defines the desired state of ThermoPilot
type ThermoPilotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// SwitchBot API credentials stored in a Secret
	// +required
	SecretRef SecretReference `json:"secretRef"`

	// Device IDs for controlling temperature
	// +optional
	AirConditionerID string `json:"airConditionerId,omitempty"`

	// Type of temperature sensor to use (e.g., MeterPro)
	// +kubebuilder:validation:Enum=MeterPro
	// +required
	TemperatureSensorType string `json:"temperatureSensorType"`

	// Temperature control settings
	// +kubebuilder:validation:Pattern=^([1-3][0-9]|[1-9])(\.[0-9])?$
	// +required
	TargetTemperature string `json:"targetTemperature"`
	// +kubebuilder:validation:Pattern=^[0-5](\.[0-9])?$
	// +kubebuilder:default="1.0"
	// +optional
	Threshold string `json:"threshold,omitempty"`

	// Air conditioner mode: cool or heat
	// +kubebuilder:validation:Enum=cool;heat
	// +required
	Mode string `json:"mode"`
}

// SecretReference holds a reference to a Secret containing SwitchBot API credentials
type SecretReference struct {
	// Name of the Secret in the same namespace
	// +required
	Name string `json:"name"`
	// Key containing the SwitchBot API token
	// +kubebuilder:default=token
	// +optional
	TokenKey string `json:"tokenKey,omitempty"`
	// Key containing the SwitchBot API secret
	// +kubebuilder:default=secret
	// +optional
	SecretKey string `json:"secretKey,omitempty"`
}

// ThermoPilotStatus defines the observed state of ThermoPilot.
type ThermoPilotStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the ThermoPilot resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// +optional
	CurrentTemperature string `json:"currentTemperature,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ThermoPilot is the Schema for the thermopilots API
type ThermoPilot struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ThermoPilot
	// +required
	Spec ThermoPilotSpec `json:"spec"`

	// status defines the observed state of ThermoPilot
	// +optional
	Status ThermoPilotStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ThermoPilotList contains a list of ThermoPilot
type ThermoPilotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ThermoPilot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ThermoPilot{}, &ThermoPilotList{})
}
