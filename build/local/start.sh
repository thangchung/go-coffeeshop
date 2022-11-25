#!/usr/bin/env bash
set -euo pipefail

require() {
  if ! hash "$1" &>/dev/null; then
    echo "'$1' not found in PATH"
    exit 1
  fi
}

require consul
require nomad

cleanup() {
  echo
  echo "Shutting down services"
  kill $(jobs -p)
  wait
}
trap cleanup EXIT

# https://github.com/deislabs/hippo/blob/de73ae52d606c0a2351f90069e96acea831281bc/src/Infrastructure/Jobs/NomadJob.cs#L28
# https://www.nomadproject.io/docs/drivers/exec#client-requirements
case "$OSTYPE" in
  linux*) SUDO="sudo --preserve-env=PATH" ;;
  *) SUDO= ;;
esac

# change to the directory of this script
cd "$(dirname "${BASH_SOURCE[0]}")"

${SUDO} rm -rf ./data
mkdir -p log

IP_ADDRESS=$(hostname -I | xargs | awk '{print $1}')

echo "Starting consul..."
consul agent -dev \
  -config-file ./etc/consul.hcl \
  -bootstrap-expect 1 \
  -client '0.0.0.0' \
  -bind "${IP_ADDRESS}" \
  &>log/consul.log &

echo "Waiting for consul..."
while ! consul members &>/dev/null; do
  sleep 2
done

echo "Starting nomad..."
${SUDO} nomad agent -dev \
  -config ./etc/nomad.hcl \
  -network-interface eth0 \
  -consul-address "${IP_ADDRESS}:8500" \
  &>log/nomad.log &

echo "Waiting for nomad..."
while ! nomad server members 2>/dev/null | grep -q alive; do
  sleep 2
done

echo
echo "Done!!!"

wait