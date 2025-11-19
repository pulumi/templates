PROJECT_NAME 	:= Pulumi Templates
TESTPARALLELISM ?= 10
TESTFLAGS       := -v -count=1 -cover -timeout 1h -parallel $(TESTPARALLELISM)


# Quick test of single template, for example:
#
#     make test.typescript
#
# Useful for verifying local changes. Uses filestate backend.
test.%:
	rm -rf ./state
	mkdir -p ./state
	cd tests && \
		PULUMI_TEMPLATE_LOCATION=${PWD} \
		PULUMI_BACKEND_URL=file://${PWD}/state \
		go test -run "TestTemplate/^$*$$" $(TESTFLAGS)


# Tests a single template, for example:
#
#    make test_template.typescript
#
# This is the more general form of `make test.typescript`.
#
# If `PULUMI_TEMPLATE_LOCATION` is unspecified, it will test the
# official template from the `pulumi/templates` repo master branch,
# that is the same location where pulumi new reads templates by
# default.
#
# To test local modifications provide `PULUMI_TEMPLATE_LOCATION=$PWD`.
#
# This target also obeys useful Pulumi environment variables:
# https://www.pulumi.com/docs/reference/cli/environment-variables/
#
# To test with the default backend provide `PULUMI_ACCESS_TOKEN=...`:
# https://www.pulumi.com/docs/intro/pulumi-service/accounts/#access-tokens
#
# To use local files as a state backend provide
# `PULUMI_BACKEND_URL=file://...`.
test_template.%:
	cd tests && go test -run "TestTemplate/^$*$$" $(TESTFLAGS)


# Every template doubles up as a benchmark, for example:
#
#     make bench_template.typescript
#
# This will run a typescript template test and populate `./traces`
# with performance data. See also
# https://www.pulumi.com/docs/support/troubleshooting/#performance
#
# Accepts the same environment variables as `test_template.%` target.
bench_template.%:
	rm -rf ./traces
	mkdir -p ./traces
	cd tests/perf && \
		PULUMI_TRACING_DIR=${PWD}/traces \
		go test -run "TestTemplatePerf/^$*$$" $(TESTFLAGS)


test_templates::
	cd tests && \
		go test $(TESTFLAGS)

metadata::
	yarn && yarn run metadata && yarn test

ensure::
	cd tests && go mod download
