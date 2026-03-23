# Node Setup Guide

This guide walks through setting up a fresh Debian machine and joining it to the k3s cluster.

It assumes:

- a clean Debian install
- SSH access to the machine
- this repo is available on the machine

## 1. Clone the repository

```bash
git clone <your-repo-url>
cd homelab
```

## 2. Run bootstrap script

The bootstrap script prepares the system and installs k3s.

### Control plane (first node)

```bash
sudo ./scripts/bootstrap.sh
```

This will:

- update the system
- install base packages
- disable swap
- configure kernel networking
- install k3s in server mode

## 3. Retrieve cluster token

After installing the first node, retrieve the node token:

```bash
sudo cat /var/lib/rancher/k3s/server/node-token
```

Save this. You will need it to join other nodes.

## 4. Add additional nodes as agents

On a new machine:

```bash
sudo K3S_ROLE=agent \
K3S_URL=https://<SERVER_IP>:6443 \
K3S_TOKEN=<TOKEN> \
./scripts/bootstrap.sh
```

Replace:

- `<SERVER_IP>` with the IP of your control plane node
- `<TOKEN>` with the token from the previous step

## 5. Verify the cluster

On the control plane node:

```bash
kubectl get nodes
```

Expected output:

```text
NAME         STATUS   ROLES                  AGE   VERSION
k3s-node-1   Ready    control-plane,master   ...   ...
k3s-node-2   Ready    <none>                 ...   ...
```

## 6. Access kubectl easily

By default, kubeconfig is stored at:

```text
/etc/rancher/k3s/k3s.yaml
```

Make it available in your shell:

```bash
echo 'export KUBECONFIG=/etc/rancher/k3s/k3s.yaml' >> ~/.bashrc
source ~/.bashrc
```

## 7. What the bootstrap script configures

### System

- installs curl, git, vim, htop, and other base packages
- disables swap
- sets the hostname

### Kernel

- enables `overlay`
- enables `br_netfilter`
- configures networking sysctl rules

### Kubernetes (k3s)

- installs k3s as a server or agent
- can disable default Traefik if configured that way
- sets kubeconfig permissions

## 8. Troubleshooting

### Node not showing up

Check the k3s service on the server node:

```bash
sudo systemctl status k3s
```

Or on an agent node:

```bash
sudo systemctl status k3s-agent
```

### Cannot connect to API server

Check:

- the `K3S_URL` is correct
- port `6443` is open
- firewall rules allow access

### Pods stuck or not scheduling

```bash
kubectl describe node <node-name>
```

## 9. Re-running bootstrap

This script is intended as a one-time setup.

If you rerun it:

- some steps may already be applied
- a clean reinstall of the OS is the safest reset path

## 10. Next steps

After nodes are set up, good next steps are:

- install monitoring
- add ingress
- set up GitOps

## 11. Notes

- this setup uses a single control plane
- it is good for learning and homelab use
- it is not a highly available production setup

## 12. Related files

```text
scripts/bootstrap.sh
scripts/install-k3s.sh
k3s/config.yaml
```
