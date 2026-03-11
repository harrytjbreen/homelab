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

## Update the system

Run:

`sudo apt update`

Then:

`sudo apt upgrade -y`

---

## Install base tooling

Run:

`sudo apt install -y git docker.io docker-compose curl htop tmux vim`

---

## Allow the user to run Docker

Run:

`sudo usermod -aG docker harry`

Log out and log back in.

Test Docker:

`docker run hello-world`

---

# Homelab Repository Setup

Clone the repository:

`git clone <repo-url> ~/homelab`

Then:

`cd ~/homelab`

Repository layout:

```
homelab/
├── apps/
├── hosts/
│   └── pi-server/
├── docs/
└── scripts/
```

---

# Deployment

Containers for this host are defined in:

`hosts/pi-server/docker-compose.yml`

Start services:

`docker compose up -d`

View running containers:

`docker ps`

---

# Updating Services

Pull the latest infrastructure configuration:

`git pull`

Redeploy containers:

`docker compose up -d`

---

# Data Storage

Persistent container data is stored in:

`/home/harry/data`

Example layout:

```
data/
├── uptime-kuma
├── grafana
└── traefik
```

---

# Rebuilding the Server

If the Pi needs to be rebuilt:

1. Flash Raspberry Pi OS Lite
2. Enable SSH
3. Install Docker
4. Clone the homelab repository
5. Run:

`docker compose up -d`

All services should automatically redeploy.

