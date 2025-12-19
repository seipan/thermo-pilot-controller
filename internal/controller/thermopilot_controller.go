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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	thermopilotv1 "github.com/seipan/thermo-pilot-controller/api/v1"
	switchbotclient "github.com/seipan/thermo-pilot-controller/internal/client"
)

// ThermoPilotReconciler reconciles a ThermoPilot object
type ThermoPilotReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=thermo-pilot.yadon3141.com,resources=thermopilots,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=thermo-pilot.yadon3141.com,resources=thermopilots/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=thermo-pilot.yadon3141.com,resources=thermopilots/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ThermoPilotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var thermoPilot thermopilotv1.ThermoPilot
	if err := r.Get(ctx, req.NamespacedName, &thermoPilot); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch ThermoPilot")
		return ctrl.Result{}, err
	}

	if thermoPilot.Status.Conditions == nil {
		thermoPilot.Status.Conditions = []metav1.Condition{}
	}

	creds, err := GetSwitchBotCredentials(ctx, r.Client, thermoPilot.Spec, thermoPilot.Namespace)
	if err != nil {
		log.Error(err, "failed to get SwitchBot credentials")
		r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "CredentialsError", err.Error())
		if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
			log.Error(statusErr, "failed to update status")
		}
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, err
	}

	sbClient := switchbotclient.NewClient(creds.Token, creds.Secret)

	// Get temperature sensor device
	var sensorID string
	switch thermoPilot.Spec.TemperatureSensorType {
	case "MeterPro":
		meterPro, err := sbClient.GetMeterPro(ctx)
		if err != nil {
			log.Error(err, "failed to get MeterPro device")
			r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "TemperatureSensorNotFound", err.Error())
			if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
				log.Error(statusErr, "failed to update status")
			}
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
		}
		sensorID = meterPro.DeviceID
		log.Info("found temperature sensor", "type", thermoPilot.Spec.TemperatureSensorType, "deviceId", sensorID, "name", meterPro.DeviceName)
	default:
		err := fmt.Errorf("unsupported temperature sensor type: %s", thermoPilot.Spec.TemperatureSensorType)
		log.Error(err, "invalid sensor type")
		r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "ConfigError", err.Error())
		if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
			log.Error(statusErr, "failed to update status")
		}
		return ctrl.Result{}, err
	}

	// Get current temperature from sensor
	currentTemp, err := sbClient.GetNowTemperature(ctx, sensorID)
	if err != nil {
		log.Error(err, "failed to get current temperature", "sensorId", sensorID)
		r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "TemperatureSensorError", err.Error())
		if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
			log.Error(statusErr, "failed to update status")
		}
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
	}

	thermoPilot.Status.CurrentTemperature = FormatTemperature(currentTemp)

	targetTemp, err := ParseTemperature(thermoPilot.Spec.TargetTemperature)
	if err != nil {
		log.Error(err, "failed to parse target temperature")
		r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "ConfigError", err.Error())
		if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
			log.Error(statusErr, "failed to update status")
		}
		return ctrl.Result{}, err
	}
	threshold := 1.0
	if thermoPilot.Spec.Threshold != "" {
		threshold, err = ParseTemperature(thermoPilot.Spec.Threshold)
		if err != nil {
			log.Error(err, "failed to parse threshold")
			threshold = 1.0
		}
	}

	tempDiff := currentTemp - targetTemp
	log.Info("temperature status",
		"current", currentTemp,
		"target", targetTemp,
		"difference", tempDiff,
		"threshold", threshold)

	needsAction := false
	var action string
	var mode switchbotclient.AirConditionerMode

	switch thermoPilot.Spec.Mode {
	case "cool":
		if tempDiff > threshold {
			needsAction = true
			action = "cooling"
			mode = switchbotclient.ModeCool
		} else if tempDiff < -threshold {
			needsAction = true
			action = "adjusting up (too cold)"
			mode = switchbotclient.ModeCool
		}
	case "heat":
		if tempDiff < -threshold {
			needsAction = true
			action = "heating"
			mode = switchbotclient.ModeHeat
		} else if tempDiff > threshold {
			needsAction = true
			action = "adjusting down (too warm)"
			mode = switchbotclient.ModeHeat
		}
	default:
		err := fmt.Errorf("unsupported mode: %s", thermoPilot.Spec.Mode)
		log.Error(err, "invalid mode in spec")
		r.setCondition(&thermoPilot, "Available", metav1.ConditionFalse, "ConfigError", err.Error())
		if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
			log.Error(statusErr, "failed to update status")
		}
		return ctrl.Result{}, err
	}

	if needsAction {
		log.Info("controlling air conditioner", "action", action, "mode", mode)
		r.setCondition(&thermoPilot, "Progressing", metav1.ConditionTrue, "ControllingAirConditioner", fmt.Sprintf("Performing action: %s", action))

		adjustedTemp := targetTemp
		if action == "adjusting up (too cold)" {
			adjustedTemp = targetTemp + 3.0 // Raise by 3 degrees if too cold
		} else if action == "adjusting down (too warm)" {
			adjustedTemp = targetTemp - 3.0 // Lower by 3 degrees if too warm
		}

		// Get air conditioner IDs
		var airConditionerIDs []string
		if thermoPilot.Spec.AirConditionerID != "" {
			airConditionerIDs = []string{thermoPilot.Spec.AirConditionerID}
		} else {
			// Get all air conditioners
			airConditioners, err := sbClient.MultiGetAirConditioners(ctx)
			if err != nil {
				log.Error(err, "failed to get air conditioners")
				r.setCondition(&thermoPilot, "Degraded", metav1.ConditionTrue, "AirConditionerListError", err.Error())
				if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
					log.Error(statusErr, "failed to update status")
				}
				return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
			}
			for _, ac := range airConditioners {
				airConditionerIDs = append(airConditionerIDs, ac.DeviceID)
			}
			log.Info("found air conditioners", "count", len(airConditionerIDs), "ids", airConditionerIDs)
		}

		// Control all air conditioners
		var controlErrors []string
		for _, deviceID := range airConditionerIDs {
			err = sbClient.SetTemperature(ctx, deviceID, adjustedTemp, mode)
			if err != nil {
				log.Error(err, "failed to control air conditioner", "deviceId", deviceID)
				controlErrors = append(controlErrors, fmt.Sprintf("%s: %v", deviceID, err))
			} else {
				log.Info("successfully controlled air conditioner", "deviceId", deviceID, "action", action)
			}
		}

		if len(controlErrors) > 0 {
			errorMsg := fmt.Sprintf("failed to control %d/%d air conditioners: %v", len(controlErrors), len(airConditionerIDs), controlErrors)
			r.setCondition(&thermoPilot, "Degraded", metav1.ConditionTrue, "AirConditionerControlError", errorMsg)
			if statusErr := r.Status().Update(ctx, &thermoPilot); statusErr != nil {
				log.Error(statusErr, "failed to update status")
			}
			if len(controlErrors) == len(airConditionerIDs) {
				return ctrl.Result{RequeueAfter: 1 * time.Minute}, fmt.Errorf(errorMsg)
			}
		}
		log.Info("air conditioner control completed", "total", len(airConditionerIDs), "errors", len(controlErrors))
	}

	if needsAction {
		r.setCondition(&thermoPilot, "Progressing", metav1.ConditionTrue, "TemperatureAdjusting",
			fmt.Sprintf("Adjusting temperature: current=%.1f, target=%.1f", currentTemp, targetTemp))
	} else {
		r.setCondition(&thermoPilot, "Progressing", metav1.ConditionFalse, "TemperatureStable",
			fmt.Sprintf("Temperature is within threshold: current=%.1f, target=%.1f", currentTemp, targetTemp))
	}

	r.setCondition(&thermoPilot, "Available", metav1.ConditionTrue, "Reconciling", "ThermoPilot is functioning normally")
	r.setCondition(&thermoPilot, "Degraded", metav1.ConditionFalse, "Healthy", "No errors detected")

	if err := r.Status().Update(ctx, &thermoPilot); err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *ThermoPilotReconciler) setCondition(thermoPilot *thermopilotv1.ThermoPilot, conditionType string, status metav1.ConditionStatus, reason, message string) {
	meta.SetStatusCondition(&thermoPilot.Status.Conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: thermoPilot.Generation,
	})
}

func (r *ThermoPilotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&thermopilotv1.ThermoPilot{}).
		Named("thermopilot").
		Complete(r)
}
