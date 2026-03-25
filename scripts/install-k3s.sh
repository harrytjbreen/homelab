#!/usr/bin/env bash
set -euo pipefail

ROLE=$1
K3S_URL=$2
K3S_TOKEN=$3

echo "Installing k3s as: $ROLE"

if [ "$ROLE" = "server" ]; then
  curl -sfL https://get.k3s.io | sh -s - \
    --write-kubeconfig-mode 644 \

elif [ "$ROLE" = "agent" ]; then
  if [ -z "$K3S_URL" ] || [ -z "$K3S_TOKEN" ]; then
    echo "K3S_URL and K3S_TOKEN required for agent"
    exit 1
  fi

  curl -sfL https://get.k3s.io | K3S_URL="$K3S_URL" K3S_TOKEN="$K3S_TOKEN" sh -
else
  echo "Unknown role: $ROLE"
  exit 1
fi

echo "k3s installation complete"
