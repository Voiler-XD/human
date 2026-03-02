# human-plugin-k8s

A Kubernetes manifest generator plugin for the Human compiler.

## What It Generates

- `namespace.yaml` — Kubernetes namespace
- `deployment.yaml` — Deployment with health probes, resource limits, and config references
- `service.yaml` — ClusterIP Service
- `ingress.yaml` — Ingress with TLS
- `configmap-<env>.yaml` — ConfigMap per environment (staging, production, etc.)

## Category

infra

## Quick Start

```bash
# Build the plugin
make build

# Install into Human
make install

# Verify installation
human plugin list
```

## Settings

Configure via `.human/config.json`:

```json
{
  "plugins": [
    {
      "name": "k8s",
      "settings": {
        "namespace": "my-custom-namespace"
      }
    }
  ]
}
```

## IR Fields Used

| Field | Usage |
|-------|-------|
| `app.Name` | Namespace, deployment, service, ingress names |
| `app.Config.Backend` | Container port selection (Go=8080, Python=8000, Node=3000) |
| `app.Config.Database` | DATABASE_URL secret reference |
| `app.Environments` | ConfigMap per environment with env vars |
| `app.Architecture` | Service topology (future: multi-deployment) |
