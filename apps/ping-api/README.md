# Ping API

A tiny Go REST API with a single `/ping` route that returns JSON.

## Run

```sh
cd /Users/harrybreen/Desktop/homelab/apps/ping-api
go run .
```

Then visit:

```sh
curl http://localhost:8080/ping
```

## Container

```sh
docker build -t ping-api:local .
```

## Kubernetes (via Argo CD)

- Update the image in `apps/ping-api/manifests/deployment.yaml` (or publish to a registry and set the full image URL).
- Ensure `platform/argocd/application.yaml` and `platform/argocd/applications/ping-api.yaml` point at your repo URL.
- Argo CD will sync `apps/ping-api/manifests` on pushes to `main`.

## Configuration

- `PING_API_ADDR`: address to listen on (default `:8080`).
