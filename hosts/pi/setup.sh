#!/usr/bin/env bash

set -euo pipefail

USER_NAME="harry"
HOME_DIR="/home/${USER_NAME}"
DATA_DIR="${HOME_DIR}/data"
HOMELAB_DIR="${HOME_DIR}/homelab"
HOST_DIR="${HOMELAB_DIR}/hosts/pi-server"

echo "==> Updating apt"
sudo apt update
sudo apt upgrade -y

echo "==> Installing base packages"
sudo apt install -y \
  git \
  curl \
  vim \
  tmux \
  htop \
  ca-certificates \
  gnupg \
  docker.io \
  docker-compose

echo "==> Enabling Docker on boot"
sudo systemctl enable docker
sudo systemctl start docker

echo "==> Adding ${USER_NAME} to docker group"
sudo usermod -aG docker "${USER_NAME}"

echo "==> Creating homelab directories"
mkdir -p "${HOMELAB_DIR}"
mkdir -p "${HOST_DIR}"
mkdir -p "${HOST_DIR}/prometheus"

echo "==> Creating data directories"
mkdir -p "${DATA_DIR}/uptime-kuma"
mkdir -p "${DATA_DIR}/prometheus"
mkdir -p "${DATA_DIR}/grafana"

echo "==> Fixing container data permissions"
sudo chown -R 65534:65534 "${DATA_DIR}/prometheus"
sudo chown -R 472:472 "${DATA_DIR}/grafana"

sudo chmod -R 755 "${DATA_DIR}/prometheus"
sudo chmod -R 755 "${DATA_DIR}/grafana"
sudo chmod -R 755 "${DATA_DIR}/uptime-kuma"

echo "==> Writing .env file if missing"
ENV_FILE="${HOST_DIR}/.env"
if [ ! -f "${ENV_FILE}" ]; then
  cat > "${ENV_FILE}" <<EOF
DATA_DIR=${DATA_DIR}
TZ=Europe/London
EOF
  echo "Created ${ENV_FILE}"
else
  echo "${ENV_FILE} already exists, skipping"
fi

echo "==> Setup complete"
echo
echo "Important:"
echo "- Docker has been enabled and will start automatically on boot."
echo "- ${USER_NAME} has been added to the docker group."
echo "- You must log out and back in before docker works without sudo."
