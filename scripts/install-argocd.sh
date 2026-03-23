#!/usr/bin/env sh
set -e

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

# Ensure the argocd namespace exists.
kubectl apply -f "$REPO_ROOT/platform/argocd/namespace.yaml"

# Install Argo CD into the cluster.
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Create the Argo CD Application for this repo.
kubectl apply -k "$REPO_ROOT/platform/argocd"
