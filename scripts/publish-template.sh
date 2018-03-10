#!/bin/bash
# publish-template.sh builds and publishes a package containing the template to
# s3://rel.pulumi.com/releases/templates.
set -o nounset -o errexit -o pipefail

# For now, all our templates have the same description.
TEMPLATE_DESCRIPTION="A Pulumi project."

ROOT=$(dirname $0)/..
TEMPLATE_SOURCE_PATH="${ROOT}/templates/$1"
TEMPLATE_PACKAGE_NAME="$1.tar.gz"
TEMPLATE_PACKAGE_DIR="$(mktemp -d)"
TEMPLATE_PACKAGE_PATH="${TEMPLATE_PACKAGE_DIR}/${TEMPLATE_PACKAGE_NAME}"

# Tar up the template
tar -czf ${TEMPLATE_PACKAGE_PATH} -C ${TEMPLATE_SOURCE_PATH} .

# rel.pulumi.com is in our production account, so assume that role first
CREDS_JSON=$(aws sts assume-role \
                 --role-arn "arn:aws:iam::058607598222:role/UploadPulumiReleases" \
                 --role-session-name "upload-plugin-pulumi-resource-aws" \
                 --external-id "upload-pulumi-release")

# Use the credentials we just assumed
export AWS_ACCESS_KEY_ID=$(echo ${CREDS_JSON}     | jq ".Credentials.AccessKeyId" --raw-output)
export AWS_SECRET_ACCESS_KEY=$(echo ${CREDS_JSON} | jq ".Credentials.SecretAccessKey" --raw-output)
export AWS_SECURITY_TOKEN=$(echo ${CREDS_JSON}    | jq ".Credentials.SessionToken" --raw-output)

# TODO
# aws s3 cp --only-show-errors "${TEMPLATE_PACKAGE_PATH}" "s3://rel.pulumi.com/releases/templates/${TEMPLATE_PACKAGE_NAME}" \
#     --metadata '{ "pulumi-template-description": "${TEMPLATE_DESCRIPTION}" }'

rm -rf "${TEMPLATE_PACKAGE_DIR}"
