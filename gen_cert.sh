#!/bin/bash

set -ue

namespace=${1:-kube-system}
svc_name=${2:-irsa-emu}
days=${3:-36500}

echo "subjectAltName = DNS:${svc_name}.${namespace}.svc" > /tmp/irsa-emu.san.txt
openssl genrsa -out irsa-emu.key.pem 2048
openssl req -new -key irsa-emu.key.pem -subj "/CN=${svc_name}.${namespace}.svc" -out /tmp/irsa-emu.csr.pem
openssl x509 -req -days "${days}" -in /tmp/irsa-emu.csr.pem -signkey irsa-emu.key.pem -out irsa-emu.cert.pem -extfile /tmp/irsa-emu.san.txt

echo
echo Base64 encoded private key:
cat irsa-emu.key.pem | base64 -w0

echo
echo
echo Base64 encoded cert:
cat irsa-emu.cert.pem | base64 -w0

echo
