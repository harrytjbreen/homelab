#!/usr/bin/env bash
set -euo pipefail

echo "Starting homelab bootstrap..."

# --- Config ---
HOSTNAME=${HOSTNAME:-k3s-node-1}
K3S_ROLE=${K3S_ROLE:-server} # server | agent
K3S_URL=${K3S_URL:-""}
K3S_TOKEN=${K3S_TOKEN:-""}

# --- Basic system setup ---
echo "Setting hostname..."
hostnamectl set-hostname "$HOSTNAME"

echo "Updating system..."
apt update && apt upgrade -y

echo "Installing base packages..."
apt install -y \
  curl \
  wget \
  git \
  vim \
  htop \
  net-tools \
  ca-certificates \
  gnupg \
  lsb-release

# --- Disable swap (required for k8s) ---
echo "Disabling swap..."
swapoff -a
sed -i '/ swap / s/^/#/' /etc/fstab

# --- Kernel modules ---
echo "Configuring kernel modules..."
cat <<EOF | tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

modprobe overlay
modprobe br_netfilter

cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sysctl --system

echo "Installing k3s..."
chmod +x ./scripts/install-k3s.sh

./scripts/install-k3s.sh "$K3S_ROLE" "$K3S_URL" "$K3S_TOKEN"

echo "Bootstrap complete!"
