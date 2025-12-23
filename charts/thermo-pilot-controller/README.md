# thermo-pilot-controller

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

A Kubernetes operator for controlling air conditioners based on temperature readings from SwitchBot sensors

## Prerequisites

- Kubernetes 1.24+
- Helm 3.0+
- SwitchBot account with API credentials

## Installation

```bash
# Add the Helm repository
helm repo add thermo-pilot https://seipan.github.io/thermo-pilot-controller
helm repo update

# Install the chart
helm install thermo-pilot thermo-pilot/thermo-pilot-controller
```

### Install with custom values

```bash
# Create values file
cat <<EOF > values.yaml
image:
  repository: your-registry/thermo-pilot-controller
  tag: v0.1.0

resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
EOF

# Install with custom values
helm install thermo-pilot thermo-pilot/thermo-pilot-controller -f values.yaml
```

## Configuration

### Required Setup

1. Create SwitchBot API credentials secret:

```bash
kubectl create secret generic switchbot-credentials \
  --from-literal=token="your-switchbot-token" \
  --from-literal=secret="your-switchbot-secret"
```

2. Create ThermoPilot resource:

```yaml
apiVersion: thermo-pilot.yadon3141.com/v1
kind: ThermoPilot
metadata:
  name: home-thermostat
spec:
  secretRef:
    name: switchbot-credentials
  temperatureSensorType: MeterPro
  targetTemperature: "25.0"
  threshold: "1.0"
  mode: cool
```

### Security Configuration

By default, the controller uses namespace-scoped access to secrets for improved security. You can configure this behavior:

```yaml
# values.yaml
rbac:
  secretAccess:
    # Use namespace-scoped access (recommended)
    namespaced: true
    # Specify a different namespace for secrets (optional)
    namespace: "switchbot-secrets"
```

For cluster-wide secret access (less secure):

```yaml
rbac:
  secretAccess:
    namespaced: false
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| replicaCount | int | `1` | Number of controller replicas |
| image.repository | string | `"ghcr.io/seipan/thermo-pilot-controller"` | Container image repository |
| image.pullPolicy | string | `"IfNotPresent"` | Container image pull policy |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.name | string | `""` | The name of the service account to use |
| podAnnotations | object | `{}` | Pod annotations |
| podSecurityContext | object | `{"runAsNonRoot":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Pod security context |
| securityContext | object | `{"allowPrivilegeEscalation":false,"readOnlyRootFilesystem":true,"capabilities":{"drop":["ALL"]}}` | Container security context |
| resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}` | Resource limits and requests |
| nodeSelector | object | `{}` | Node selector |
| tolerations | list | `[]` | Tolerations |
| affinity | object | `{}` | Affinity rules |
| podDisruptionBudget.enabled | bool | `false` | Enable PodDisruptionBudget |
| podDisruptionBudget.minAvailable | int | `1` | Minimum available pods during disruption |
| podDisruptionBudget.maxUnavailable | string | `""` | Maximum unavailable pods during disruption |
| controller.leaderElect | bool | `true` | Enable leader election |
| controller.healthProbeBindAddress | string | `":8081"` | Health probe bind address |
| controller.metricsBindAddress | string | `":8080"` | Metrics bind address |
| controller.metricsSecure | bool | `true` | Enable secure metrics endpoint |
| probes.liveness.enabled | bool | `true` | Enable liveness probe |
| probes.liveness.initialDelaySeconds | int | `15` | Initial delay seconds |
| probes.liveness.periodSeconds | int | `20` | Period seconds |
| probes.readiness.enabled | bool | `true` | Enable readiness probe |
| probes.readiness.initialDelaySeconds | int | `5` | Initial delay seconds |
| probes.readiness.periodSeconds | int | `10` | Period seconds |
| rbac.create | bool | `true` | Create RBAC resources |
| rbac.secretAccess.namespaced | bool | `true` | Use namespace-scoped secret access |
| rbac.secretAccess.namespace | string | `""` | Namespace where secrets can be accessed |
| monitoring.enabled | bool | `false` | Enable ServiceMonitor for Prometheus Operator |
| monitoring.labels | object | `{}` | ServiceMonitor labels |
| monitoring.interval | string | `"30s"` | ServiceMonitor scrape interval |
| thermoPilot.enabled | bool | `false` | Enable example ThermoPilot resource |

## Uninstallation

```bash
# Remove ThermoPilot resources first
kubectl delete thermopilots --all

# Uninstall the chart
helm uninstall thermo-pilot

# Remove CRDs (optional)
kubectl delete crd thermopilots.thermo-pilot.yadon3141.com
```

## Support

For issues and feature requests, please visit:
https://github.com/seipan/thermo-pilot-controller/issues