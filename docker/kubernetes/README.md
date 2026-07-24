# Kubernetes deployment

These manifests target the local Kubernetes cluster provided by Docker Desktop
while retaining production-style resource limits, probes, rolling updates, and
autoscaling.

## Prerequisites

- The main backend deployment must create the shared `bank` namespace.
- Build the application image as `bank-exchange-rate-service:local`.
- The main backend deployment must provide the shared Kafka broker through a
  Service named `kafka` on port `9092`.
- Enable Kubernetes Metrics Server for the HPA.
- Replace the local Redis password in `01-config.yaml` before non-local use.

## Deploy

```powershell
kubectl apply -f docker/kubernetes
kubectl get pods -n bank
kubectl get services -n bank
kubectl get hpa -n bank
```

The gRPC endpoint is available only inside the cluster at:

```text
exchange-rate-service.bank.svc.cluster.local:55053
```

Redis is available only inside the cluster at:

```text
exchange-rate-redis.bank.svc.cluster.local:16379
```

The application expects the Kafka broker owned by the main backend deployment
at:

```text
kafka.bank.svc.cluster.local:9092
```

The namespace and Kafka manifests intentionally do not belong to this
microservice. They must be managed from the main backend or a shared
infrastructure repository.
