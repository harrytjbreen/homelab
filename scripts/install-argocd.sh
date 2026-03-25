#!/usr/bin/env sh
set -e

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)"

# Clean up any Argo CD resources accidentally installed in the default namespace.
ARGOCD_INSTALL_URL="https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml"
if kubectl -n default get deploy argocd-server >/dev/null 2>&1; then
  kubectl delete -n default -f "$ARGOCD_INSTALL_URL" --ignore-not-found
fi

# Ensure the argocd namespace exists.
kubectl apply -f "$REPO_ROOT/platform/argocd/namespace.yaml"

# Install Argo CD into the cluster (server-side avoids oversized last-applied annotations).
kubectl apply --server-side --force-conflicts --field-manager=argocd-bootstrap -n argocd \
  -f "$ARGOCD_INSTALL_URL"

# Create the Argo CD Application for this repo.
kubectl apply --server-side --field-manager=argocd-bootstrap -k "$REPO_ROOT/platform/argocd"
