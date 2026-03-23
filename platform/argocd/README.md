# Argo CD setup

This directory bootstraps Argo CD and defines a root app-of-apps that manages child applications.

## What to edit
- `application.yaml`: set `spec.source.repoURL` to your Git repository URL.
- `platform/argocd/applications/ping-api.yaml`: set `spec.source.repoURL` to your Git repository URL.

## Apply
1. Install Argo CD in the cluster (see the script in `scripts/install-argocd.sh`).
2. Apply this kustomization to create the Argo CD namespace and root Application.

Child apps live in `platform/argocd/applications` and will auto-sync on pushes to `main`.
