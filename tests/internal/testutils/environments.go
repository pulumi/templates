package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	ptesting "github.com/pulumi/pulumi/sdk/v3/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

// Allocates a new testing environment and schedules its deletion on
// test cleanup.
func NewEnvironment(t *testing.T, cfg TemplateTestConfig) *ptesting.Environment {
	e := ptesting.NewEnvironment(t)
	t.Cleanup(func() { deleteIfNotFailed(e, cfg) })
	return e
}

// deleteIfNotFailed deletes the files in the testing environment if the testcase has
// not failed. (Otherwise they are left to aid debugging.)
func deleteIfNotFailed(e *ptesting.Environment, cfg TemplateTestConfig) {
	if !e.T.Failed() {
		e.DeleteEnvironment()
	}
}

// Calls pulumi new with a given template.
//
// Since pulumi new expects a stack name or assumes dev, we generate a
// random one here to prevent conflicts. Note that ProgramTest will
// use its own stack, so we take care to delete this one right away.
//
// There is a --generate-only option that opts out of installing
// dependencies, but we do want that to happen as part of pulumi new.
func PulumiNew(e *ptesting.Environment, templatePath string, extraArgs ...string) {
	tempStack := (&integration.ProgramTestOptions{}).GetStackName().String()
	cmdArgs := append(
		[]string{"new", templatePath, "-f", "--yes", "-s", tempStack},
		extraArgs...,
	)
	e.RunCommand("pulumi", cmdArgs...)
	e.RunCommand("pulumi", "stack", "rm", tempStack, "--yes")
}

// Overrides PrepareProject ProgramTest options by auto-detecting
// environment runtime.
//
// Default PrepareProject for Node uses yarn install to install
// dependencies; template tests do not need it because pulumi new
// already installs them with npm, which is also what will happen on
// user systems.
func PrepareProject(t *testing.T, e *ptesting.Environment) func(*engine.Projinfo) error {
	path, err := workspace.DetectProjectPathFrom(e.RootPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	projinfo, err := workspace.LoadProject(path)
	assert.NoError(t, err)

	var prepareProject func(*engine.Projinfo) error
	switch rt := projinfo.Runtime.Name(); rt {
	case integration.NodeJSRuntime:
		prepareProject = func(*engine.Projinfo) error {
			return nil
		}
	default:
		prepareProject = nil // use default logic
	}
	return prepareProject
}
