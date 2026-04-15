# Argo CD setup

This directory bootstraps Argo CD and defines a root app-of-apps that manages child applications.

## What to edit
- `applications/homelab.yaml`: set `spec.source.repoURL` to your Git repository URL.
- `platform/argocd/applications/ping-api.yaml`: set `spec.source.repoURL` to your Git repository URL.

## Apply
1. Install Argo CD in the cluster (see the script in `scripts/install-argocd.sh`).
2. Apply this kustomization to create the Argo CD namespace and root Application.

Note: Argo CD should only run in the `argocd` namespace. The install script will remove any Argo CD resources that were accidentally created in `default`.

Child apps live in `platform/argocd/applications` and will auto-sync on pushes to `main`.

`argocd-image-updater` runs in the cluster and watches the `latest` tags for app images. It uses the digest update strategy so a new image pushed to the same `latest` tag still changes the rendered Deployment and rolls the pods.
