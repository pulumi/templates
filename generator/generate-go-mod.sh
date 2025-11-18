#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

GO_VERSION="1.24"

# Fetch the latest release version of the given Pulumi repo.
fetch_latest_version() {
	local repo="pulumi/$1"
	curl -s \
		-H "Authorization: token ${GITHUB_TOKEN}" \
		"https://api.github.com/repos/${repo}/releases/latest" | jq -r .name
}

PULUMI_VERSION="$(fetch_latest_version pulumi)"

PROVIDER_LIST="aiven,alicloud,auth0,aws,azure,azure-classic,civo,digitalocean,equinix-metal,gcp,google-native,kubernetes,linode,oci,openstack,random"
IFS=',' read -ra PROVIDERS <<<"$PROVIDER_LIST"

for i in "${PROVIDERS[@]}"; do
	echo "Updating $i template"
	if [ "$i" = "azure-classic" ]; then
		PROVIDER_NAME="azure"
	elif [ "$i" = "azure" ]; then
		PROVIDER_NAME="azure-native"
	else
		PROVIDER_NAME="$i"
	fi
	PROVIDER_VERSION="$(fetch_latest_version "pulumi-$PROVIDER_NAME")"
	MAJOR_VERSION=""
	if [[ "$PROVIDER_VERSION" =~ ^v([0-9]+)\. ]]; then
		major_num="${BASH_REMATCH[1]}"
		if [ "$major_num" -gt 1 ]; then
			MAJOR_VERSION="/v$major_num"
		fi
	fi
	REQUIRE_LINE="github.com/pulumi/pulumi-${PROVIDER_NAME}/sdk${MAJOR_VERSION} ${PROVIDER_VERSION}"
	if [ "$PROVIDER_NAME" = "azure-native" ]; then
		# azure-native has a different internal structure. Generate the go.mod file accordingly.
		REQUIRE_LINE="github.com/pulumi/pulumi-azure-native-sdk/resources${MAJOR_VERSION} ${PROVIDER_VERSION}
	github.com/pulumi/pulumi-azure-native-sdk/storage${MAJOR_VERSION} ${PROVIDER_VERSION}"
	fi

        cat<<EOF > "../$i-go/go.mod"
module \${PROJECT}

go ${GO_VERSION}

require (
	${REQUIRE_LINE}
	github.com/pulumi/pulumi/sdk/v3 ${PULUMI_VERSION}
)
EOF
	echo "Updated $i go mod template to be Pulumi $PULUMI_VERSION and pulumi-$i to be $PROVIDER_VERSION"
done

echo "Updating go template"
cat<<EOF > ../go/go.mod
module \${PROJECT}

go ${GO_VERSION}

require (
	github.com/pulumi/pulumi/sdk/v3 ${PULUMI_VERSION}
)
EOF
echo "Updated go template to be Pulumi $PULUMI_VERSION"
