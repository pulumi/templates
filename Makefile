PROJECT_NAME := Pulumi Templates
include _tools/build/common.mk

test_templates::
	cd tests && go test -v -count=1 -cover -timeout 1h -parallel 10 .

ensure::
	cd tests && GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor

.PHONY: publish
publish:
	$(call STEP_MESSAGE)
	./_tools/scripts/publish.sh
