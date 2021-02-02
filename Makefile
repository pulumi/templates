PROJECT_NAME := Pulumi Templates

test_templates::
	cd tests && go test -v -count=1 -cover -timeout 6h -parallel 10 .

ensure::
	cd tests && go mod tidy && go mod download
