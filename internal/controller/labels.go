package controller

import (
	thermopilotv1 "github.com/seipan/thermo-pilot-controller/api/v1"
)

const (
	appName     = "thermopilot"
	managerName = "thermo-pilot-controller"
)

// appLabels returns the standard labels for ThermoPilot resources
func appLabels(thermoPilot thermopilotv1.ThermoPilot) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       appName,
		"app.kubernetes.io/instance":   thermoPilot.Name,
		"app.kubernetes.io/created-by": managerName,
		"app.kubernetes.io/part-of":    "thermo-pilot-system",
	}
}