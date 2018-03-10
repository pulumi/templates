#!/bin/bash
# publish.sh builds and publishes a release.
set -o nounset -o errexit -o pipefail

ROOT=$(dirname $0)/..

echo "Publishing templates to s3://rel.pulumi.com/:"
for template in "${ROOT}/templates"/*
do
  name=$(basename "$template")
  ${ROOT}/scripts/publish-template.sh "$name"
done
