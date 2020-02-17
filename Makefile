PROJECT_NAME := Pulumi Templates

test_templates::
	cd tests && go test -v -count=1 -cover -timeout 1h -parallel 10 .

ensure::
	cd tests && GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor
