#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

GO_VERSION="1.26"

# Fetch the latest release version of the given Pulumi repo.
fetch_latest_version() {
	local repo="pulumi/$1"
	curl -s \
		-H "Authorization: token ${GITHUB_TOKEN}" \
		"https://api.github.com/repos/${repo}/releases/latest" | jq -r .name
}

PULUMI_VERSION="$(fetch_latest_version pulumi)"

PROVIDER_LIST="aiven,alicloud,auth0,aws,azure,digitalocean,gcp,kubernetes,linode,oci,random"
IFS=',' read -ra PROVIDERS <<<"$PROVIDER_LIST"

for i in "${PROVIDERS[@]}"; do
	echo "Updating $i template"
        if [ "$i" = "azure" ]; then
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

# Other Go templates keep their own requires: bump each github.com/pulumi module within its major version, so import paths stay valid.
fetch_latest_module_version() {
	curl -sf "https://proxy.golang.org/$1/@latest" | jq -r .Version
}

set_module_version() {
	local mod_file="$1" module="$2" version="$3" tmp
	tmp="$(mktemp)"
	awk -v module="$module" -v version="$version" '
		$1 == module && $3 != "//" { sub(/[^ \t]+$/, version) }
		$1 == "require" && $2 == module && $4 != "//" { $3 = version }
		{ print }
	' "$mod_file" >"$tmp" && mv "$tmp" "$mod_file"
}

GENERATED_TEMPLATES=("${PROVIDERS[@]/%/-go}" "go")

for mod_file in ../*-go/go.mod; do
	template="$(basename "$(dirname "$mod_file")")"
	if [[ " ${GENERATED_TEMPLATES[*]} " == *" ${template} "* ]]; then
		continue
	fi
	modules="$(awk '
		$1 ~ /^github\.com\/pulumi\// && $3 != "//" { print $1 }
		$1 == "require" && $2 ~ /^github\.com\/pulumi\// && $4 != "//" { print $2 }
	' "$mod_file")"
	for module in $modules; do
		version="$(fetch_latest_module_version "$module")"
		if [ -z "$version" ] || [ "$version" = "null" ]; then
			echo "Skipping $module in $template: no version from the module proxy"
			continue
		fi
		set_module_version "$mod_file" "$module" "$version"
	done
	echo "Updated the Pulumi modules in $template"
done

