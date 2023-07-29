#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

# Fetch the latest release version of the given Pulumi repo.
fetch_latest_version() {
    local repo="pulumi/$1"
    curl -s \
         -H "Authorization: token ${GITHUB_TOKEN}" \
         "https://api.github.com/repos/${repo}/releases/latest" | jq -r .name
}

PULUMI_VERSION="$(fetch_latest_version pulumi)"

PROVIDER_LIST="aiven,alicloud,auth0,aws,aws-native,azure,azure-classic,civo,digitalocean,equinix-metal,gcp,google-native,kubernetes,linode,oci,openstack"
IFS=',' read -ra PROVIDERS <<< "$PROVIDER_LIST"

for i in "${PROVIDERS[@]}"
do
  echo "Updating" $i "template"
  if [ "$i" = "azure-classic" ]; then
    PROVIDER_NAME="azure"
  elif [ "$i" = "azure" ]; then
    PROVIDER_NAME="azure-native"
  else
    PROVIDER_NAME="$i"
  fi
  PROVIDER_VERSION="$(fetch_latest_version "pulumi-$PROVIDER_NAME")"
  sed -e "s/\${VERSION}/$PROVIDER_VERSION/g" -e "s/\${PULUMI_VERSION}/$PULUMI_VERSION/g" mod-templates/$i-template.txt | tee ../$i-go/go.mod
  echo "Updated $i go mod template to be Pulumi $PULUMI_VERSION and pulumi-$i to be" $PROVIDER_VERSION
done

echo "Updating go template"
sed -e "s/\${PULUMI_VERSION}/$PULUMI_VERSION/g"  mod-templates/go-template.txt | tee ../go/go.mod
echo "Updated go template to be Pulumi $PULUMI_VERSION"
