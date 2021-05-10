PROJECT_NAME 	:= Pulumi Templates
TESTPARALLELISM ?= 10

test_templates::
	cd tests && go test -v -count=1 -cover -timeout 6h -parallel $(TESTPARALLELISM) .

ensure::
	cd tests && go mod download