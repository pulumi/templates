#!/bin/bash
# publish.sh builds and publishes a release.
set -o nounset -o errexit -o pipefail

ROOT=$(dirname $0)/../..

echo "Publishing templates to s3://rel.pulumi.com/:"
for template in "${ROOT}"/*
do
  # Skip files (we only care about directories).
  if [ -f "$template" ]; then
    continue
  fi

  name=$(basename "$template")

  # Skip the _tools directory.
  if [ $name == '_tools' ]; then
    continue
  fi

  ${ROOT}/_tools/scripts/publish-template.sh "$name"
done
