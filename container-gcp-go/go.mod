module ${PROJECT}

go 1.21.7

toolchain go1.22.1

require (
	github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild v0.0.3
	github.com/pulumi/pulumi-gcp/sdk/v7 v7.0.0
	github.com/pulumi/pulumi-random/sdk/v4 v4.13.2
	github.com/pulumi/pulumi/sdk/v3 v3.128.0
)
