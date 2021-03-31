#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ${GOPATH}/pkg/mod/k8s.io/code-generator@v0.20.5)}

echo "脚本目录:$SCRIPT_ROOT"
echo "path :$CODEGEN_PKG"

bash "${CODEGEN_PKG}/generate-groups.sh" all \
  "github.com/cuijxin/k8s-dashboard/src/backend/plugin/client" "github.com/cuijxin/k8s-dashboard/src/backend/plugin/apis"  \
  "apis:v1alpha1" \
  --go-header-file "${SCRIPT_ROOT}"/codegen/boilerplate.go.txt
