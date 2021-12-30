#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

PULUMI_VERSION="$(curl -s https://api.github.com/repos/pulumi/pulumi/releases/latest | jq -r .name)"

PROVIDER_LIST="alicloud,aws,azure-classic,azure,digitalocean,equinix-metal,gcp,google-native,kubernetes,linode,openstack,civo,aiven,auth0,github"
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
  PROVIDER_VERSION="$(curl -s https://api.github.com/repos/pulumi/pulumi-$PROVIDER_NAME/releases/latest | jq -r .name)"
  sed -e "s/\${VERSION}/$PROVIDER_VERSION/g" -e "s/\${PULUMI_VERSION}/$PULUMI_VERSION/g" mod-templates/$i-template.txt | tee ../$i-go/go.mod
  echo "Updated $i go mod template to be Pulumi $PULUMI_VERSION and pulumi-$i to be" $PROVIDER_VERSION
done

echo "Updating go template"
sed -e "s/\${PULUMI_VERSION}/$PULUMI_VERSION/g"  mod-templates/go-template.txt | tee ../go/go.mod
echo "Updated go template to be Pulumi $PULUMI_VERSION"
