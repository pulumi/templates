PROJECT_NAME 	:= Pulumi Templates
TESTPARALLELISM ?= 10
TESTFLAGS       := -v -count=1 -cover -timeout 6h -parallel $(TESTPARALLELISM)

test_templates::
	cd tests && go test $(TESTFLAGS)

# Run a test of a single template.
# Example: make test_template.typescript
# This will run a test corresponding to the typescript template
test_template.%:
	cd tests && BLACK_LISTED_TESTS=none go test -run "TestTemplate/^$*$$" $(TESTFLAGS)

ensure::
	cd tests && go mod download
