#!/usr/bin/env bash
set -euo pipefail

# install apt-add-repository
apt update 
apt install software-properties-common -y
apt update

echo "Adding HashiCorp GPG key and repo..."
curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -
apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
apt-get update

# install cni plugins https://www.nomadproject.io/docs/integrations/consul-connect#cni-plugins
echo "Installing cni plugins..."
curl -L -o cni-plugins.tgz "https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-$( [ $(uname -m) = aarch64 ] && echo arm64 || echo amd64)"-v1.1.1.tgz
sudo mkdir -p /opt/cni/bin
sudo tar -C /opt/cni/bin -xzf cni-plugins.tgz
sudo rm ./cni-plugins.tgz

echo "Installing Consul..."
sudo apt-get install consul -y

echo "Installing Nomad..."
sudo apt-get install nomad -y

echo "Installing Vault..."
sudo apt-get install vault -y

# # configuring environment
# sudo -H -u root nomad -autocomplete-install
# sudo -H -u root consul -autocomplete-install
# sudo -H -u root vault -autocomplete-install
# sudo tee -a /etc/environment <<EOF
# export VAULT_ADDR=http://localhost:8200
# export VAULT_TOKEN=root
# EOF

source /etc/environment

# WSL2-hack - Nomad cannot run on wsl2 image, then we need to work-around
sudo mkdir -p /lib/modules/$(uname -r)/
echo '_/bridge.ko' | sudo tee -a /lib/modules/$(uname -r)/modules.builtin