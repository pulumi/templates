package main

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks int

// Mock calls to create new resources and return a canned response.
func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	// Here, we're returning a same-shaped object for all resource types.
	// We could, however, use the arguments passed into this function to
	// customize the mocked-out properties of a particular resource.
	// See the unit-testing docs for details:
	// https://www.pulumi.com/docs/iac/concepts/testing/unit/
	return args.Name + "_id", args.Inputs, nil
}

// Mock function calls and return an empty response.
func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return resource.PropertyMap{}, nil
}

func TestInfrastructure(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		var wg sync.WaitGroup

		// Example test. To run, uncomment and run `go test`.
		// infra, err := createInfrastructure(ctx)
		// assert.NoError(t, err)

		// wg.Add(1)
		// infra.bucket.Tags.ApplyT(func(tags map[string]string) error {
		// 	assert.NotNil(t, tags)
		// 	assert.Contains(t, tags, "Name")

		// 	wg.Done()
		// 	return nil
		// })

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}
