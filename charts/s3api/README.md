# s3Api Helm Chart

## Chart Details
This chart will do the following by default:

* 3 x instances with port 8080 exposed on an external LoadBalancer or ingress
* All using Kubernetes Deployments

## Configuration

The following tables lists the configurable parameters of the Spark chart and their default values.

### s3Api

|       Parameter           |           Description               |                         Default                          |
|---------------------------|-------------------------------------|----------------------------------------------------------|
| `name`                    | k8s selector key                    | `s3api`                                                  |
| `image.name`              | app name                            | `s3api`                                                  |
| `image.tag`               | Container image tag                 | `latest`                                                 |
| `image.pullPolicy`        | Container pull policy               | `Always`                                                 |
| `image.replicaCount`      | k8s deployment replicas             | `3`                                                      |
| `image.resources`         | container requested cpu             | `{cpu: "100m", memory: "128Mi"}`                         |
| `image.pullSecrets`       | image pull secret for private repo  |                                                          |
| `livenessProbe.enabled`   | Enable Liveness probe               | `true`                                                   |
| `livenessProbe.path`      | Liveness probe http path            | `/healthz`                                               |
| `readinessProbe.enabled`  | Enable Readiness probe              | `true`                                                   |
| `readinessProbe.path`     | Readiness probe http path           | `/healthz`                                               |
| `service.type`            | k8s service type                    | `ClusterIP`                                              |
| `service.internalPort`    | Container listening port            | `8888`                                                   |
| `service.externalPort`    | k8s service port                    | `80`                                                     |
| `ingress.enabled`         | Enable ingress controller           | `false`                                                  |
| `ingress.hosts`           | Ingress hosts                       | `[]`                                                     |
| `ingress.annotations`     | Ingress annotations                 | `kubernetes.io/ingress.class: nginx`                     |
| `ingress.tls`             | Ingress tls enabled                 | `false`                                                  |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml --verify levo-charts/croc-hunter-0.2.0.tgz
```

> **Tip**: You can use the default [values.yaml](values.yaml)
