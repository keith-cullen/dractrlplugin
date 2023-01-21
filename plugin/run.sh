#!/bin/bash

KUBECONFIG="/home/${USER}/.kube/config" \
KUBERNETES_MASTER="http://localhost:8080/" \
KUBERNETES_SERVICE_HOST=$(hostname) \
KUBERNETES_SERVICE_PORT="6443" \
./bin/draplugin -f "/home/${USER}/.kube/config"
