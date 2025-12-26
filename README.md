# thermo-pilot-controller

A Kubernetes operator that automatically manages room temperature using SwitchBot smart home devices.

## Overview

The thermo-pilot-controller monitors temperature sensors and controls air conditioners to maintain target temperatures. It integrates with SwitchBot API to read temperature data and send commands to air conditioning units.

## Features

- 🌡️ Automatic temperature control based on target and threshold settings
- ❄️ Support for both cooling and heating modes
- 🔄 Continuous monitoring with 5-minute reconciliation intervals
- 🏠 Multi-AC support - controls all discovered air conditioners
- 🔐 Secure credential management using Kubernetes Secrets
- 📊 Detailed status reporting with condition tracking

## Prerequisites

- Kubernetes cluster (v1.28+)
- SwitchBot account with API credentials
- Compatible devices:
  - SwitchBot MeterPro (temperature sensor)
  - SwitchBot Hub + IR-controlled air conditioner

## Installation

### Using Helm

```bash
helm repo add thermo-pilot https://seipan.github.io/thermo-pilot-controller
helm install thermo-pilot thermo-pilot/thermo-pilot-controller
```

## Usage

### 1. Create SwitchBot Credentials

First, create a Secret with your SwitchBot API credentials:

```bash
kubectl create secret generic switchbot-credentials \
  --from-literal=token="your-switchbot-token" \
  --from-literal=secret="your-switchbot-secret" \
  -n default
```

Or using YAML:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: switchbot-credentials
  namespace: default
type: Opaque
stringData:
  token: "your-switchbot-token"
  secret: "your-switchbot-secret"
```

### 2. Create a ThermoPilot Resource

Create a ThermoPilot custom resource to start temperature control:

```yaml
apiVersion: thermo-pilot.yadon3141.com/v1
kind: ThermoPilot
metadata:
  name: living-room
  namespace: default
spec:
  # API credentials reference
  secretRef:
    name: switchbot-credentials
  
  # Temperature settings
  targetTemperature: "22.0"  # Target: 22°C
  threshold: "1.0"           # ±1°C tolerance
  mode: cool                 # cool or heat
  
  # Device configuration
  temperatureSensorType: MeterPro
  # airConditionerId: "optional-device-id"  # Omit to control all ACs
```

### 3. Check Status

Monitor the temperature control status:

```bash
kubectl get thermopilot living-room -o yaml
```

Example status output:
```yaml
status:
  currentTemperature: "23.5"
  conditions:
  - type: Available
    status: "True"
    reason: Reconciling
    message: ThermoPilot is functioning normally
  - type: Progressing
    status: "True"
    reason: TemperatureAdjusting
    message: "Adjusting temperature: current=23.5, target=22.0"
```

## Configuration

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `secretRef.name` | Name of the Secret containing SwitchBot credentials | Yes | - |
| `secretRef.tokenKey` | Key for API token in the Secret | No | `token` |
| `secretRef.secretKey` | Key for API secret in the Secret | No | `secret` |
| `targetTemperature` | Desired temperature (1.0-39.0°C) | Yes | - |
| `threshold` | Temperature tolerance (0.0-5.0°C) | No | `1.0` |
| `mode` | Operating mode (`cool` or `heat`) | Yes | - |
| `temperatureSensorType` | Type of temperature sensor | Yes | `MeterPro` |
| `airConditionerId` | Specific AC device ID | No | All ACs |

## How It Works

1. **Temperature Monitoring**: Reads current temperature from SwitchBot MeterPro every 5 minutes
2. **Decision Making**: 
   - Cool mode: Activates cooling if temperature > target + threshold
   - Heat mode: Activates heating if temperature < target - threshold
3. **Smart Control**: Adjusts AC temperature by ±3°C from target when corrections are needed
4. **Status Updates**: Reports current temperature and control actions via Kubernetes status


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Kubebuilder](https://kubebuilder.io/)
- Integrates with [SwitchBot API](https://github.com/OpenWonderLabs/SwitchBotAPI)
