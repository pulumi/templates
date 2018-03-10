PROJECT_NAME := Pulumi Templates
include build/common.mk

.PHONY: publish
publish:
	$(call STEP_MESSAGE)
	./scripts/publish.sh

# The travis_* targets are entrypoints for CI.
.PHONY: travis_cron travis_push travis_pull_request travis_api
travis_cron: all
travis_push: only_build only_test
travis_pull_request: all
travis_api: all
