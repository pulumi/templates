PROJECT_NAME := Pulumi Templates
include _tools/build/common.mk

TESTPARALLELISM := 10

test_templates::
	cd tests && $(GO_TEST) -v .

ensure::
	cd tests && GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor

.PHONY: publish
publish:
	$(call STEP_MESSAGE)
	./_tools/scripts/publish.sh
