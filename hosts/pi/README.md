# Raspberry Pi Server

## Overview

This Raspberry Pi acts as a **homelab infrastructure node**.

It runs containerised services using Docker and is managed through the `homelab` repository.

Typical responsibilities include:

- monitoring
- automation
- reverse proxy
- small infrastructure services

The goal is that this machine can be **fully rebuildable using only this repository and this document**.

---

# Hardware

| Component | Value |
|---|---|
| Device | Raspberry Pi 3B+ |
| CPU | Quad-core ARM Cortex-A53 |
| RAM | 1GB |
| Storage | SD Card |
| Network | Ethernet |
| OS | Raspberry Pi OS Lite (Debian based) |

---

# Network

| Item | Value |
|---|---|
| Hostname | pi-server |
| SSH User | harry |

SSH access:

`ssh harry@<ip-address>`

Example:

`ssh harry@192.168.1.50`

A **DHCP reservation should be configured on the router** so the IP address remains consistent.

---

# Setup

## Flash Raspberry Pi OS

1. Download **Raspberry Pi Imager**
2. Insert the SD card into your computer
3. Select **Raspberry Pi OS Lite (64-bit)**

Before flashing, configure:

- hostname: `pi-server`
- username: `harry`
- enable SSH

Flash the OS to the SD card.

---

## First Boot

1. Insert the SD card into the Raspberry Pi
2. Connect ethernet
3. Connect power
4. Wait ~60 seconds for the system to boot

---

## SSH Into The Pi

From your computer:

`ssh harry@raspberrypi.local`

If that does not resolve, find the IP address from your router and run:

`ssh harry@<ip-address>`

Example:

`ssh harry@192.168.1.50`

---

# System Configuration

## Clone the repository

Clone the homelab repository:

`git clone <repo-url> ~/homelab`

Then:

`cd ~/homelab`

Repository layout:

```text
homelab/
├── apps/
├── hosts/
│   └── pi/
├── docs/
└── scripts/
```

---

## Run the setup script

The initial machine bootstrap is handled by the host setup script.

Move into the host directory:

`cd ~/homelab/hosts/pi`

Make the script executable:

`chmod +x setup.sh`

Run it:

`./setup.sh`

This script is responsible for:

- updating the system
- installing base packages
- installing Docker and Docker Compose
- enabling Docker to start on boot
- creating required data directories
- setting required permissions for container volumes
- writing the host `.env` file if needed

---

## Re-login after setup

The setup script adds the `harry` user to the `docker` group.

Log out and log back in before using Docker without `sudo`.

Test Docker:

`docker run hello-world`

---

# Deployment

Containers for this host are defined in:

`hosts/pi/docker-compose.yml`

Move into the host directory:

`cd ~/homelab/hosts/pi`

Start services:

`docker compose up -d`

View running containers:

`docker ps`

---

# Updating Services

Pull the latest infrastructure configuration:

`git pull`

Then from `hosts/pi` redeploy containers:

`docker compose up -d`

---

# Data Storage

Persistent container data is stored in:

`/home/harry/data`

Example layout:

```text
data/
├── uptime-kuma
├── prometheus
└── grafana
```

---

# Rebuilding the Server

If the Pi needs to be rebuilt:

1. Flash Raspberry Pi OS Lite
2. Enable SSH
3. Clone the homelab repository
4. Move to `~/homelab/hosts/pi`
5. Run `chmod +x setup.sh`
6. Run `./setup.sh`
7. Log out and back in
8. Run `docker compose up -d`

All services should automatically redeploy.

---

# Planned Services

The following services are expected to run on this node:

- uptime-kuma (monitoring)
- prometheus (metrics collection)
- grafana (dashboards)
- node-exporter (host metrics)

Additional services will be added as the homelab grows.

---

# Future Improvements

- migrate OS to USB SSD
- assign static IP
- automated deployment from GitHub
- centralised monitoring across all hosts
