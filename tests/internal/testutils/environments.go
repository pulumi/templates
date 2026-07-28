package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

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
	case integration.GoRuntime:
		prepareProject = prepareGoProject
	default:
		prepareProject = nil // use default logic
	}
	return prepareProject
}

// Resolves the sdk/v3 master branch to a concrete pseudo-version once per
// test process. Resolving the branch in every test hits the module proxy
// ~20 times in parallel and intermittently fails on fresh, uncached master
// commits; `go get` then silently falls back to the repo root module and
// the install fails (see pulumi/templates#1166).
var resolveDevSDKVersion = sync.OnceValues(func() (string, error) {
	dir, err := os.MkdirTemp("", "sdk-version-resolve")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	goMod := "module resolve\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o600); err != nil {
		return "", err
	}

	var lastErr error
	for attempt := range 5 {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 10 * time.Second)
		}
		cmd := exec.Command("go", "list", "-m", "-f", "{{.Version}}", "github.com/pulumi/pulumi/sdk/v3@master")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			lastErr = fmt.Errorf("resolving sdk/v3@master: %w: %s", err, out)
			continue
		}
		return strings.TrimSpace(string(out)), nil
	}
	return "", lastErr
})

// Replaces the default Go project preparation from ProgramTest, which runs
// `go get -u .../sdk/v3@master` per test. Installs the dev SDK at the
// pinned version resolved by resolveDevSDKVersion instead.
func prepareGoProject(projinfo *engine.Projinfo) error {
	cwd, _, err := projinfo.GetPwdMain()
	if err != nil {
		return err
	}

	version, err := resolveDevSDKVersion()
	if err != nil {
		return err
	}

	steps := [][]string{
		{"go", "get", "-u", "github.com/pulumi/pulumi/sdk/v3@" + version},
		{"go", "mod", "tidy"},
	}
	for _, step := range steps {
		cmd := exec.Command(step[0], step[1:]...)
		cmd.Dir = cwd
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %w: %s", strings.Join(step, " "), err, out)
		}
	}
	return nil
}
