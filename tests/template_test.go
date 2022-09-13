package tests

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	ptesting "github.com/pulumi/pulumi/sdk/v3/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

const testTimeout = 60 * time.Minute

func TestTemplates(t *testing.T) {
	cfg := newTemplateTestConfigFromEnv()

	// When tracing is enabled to collect performance data, using
	// Quick: true skews the measurements, therefore prefer Quick:
	// false in that case.
	quick := !isTracingEnabled()

	base := integration.ProgramTestOptions{
		ExpectRefreshChanges:   true,
		Quick:                  quick,
		SkipRefresh:            true,
		NoParallel:             true, // we mark tests as Parallel manually when instantiating
		DestroyOnCleanup:       true,
		UseAutomaticVirtualEnv: true,
	}

	for _, templateInfo := range findAllTemplates(t, cfg.templateUrl) {
		templateInfo := templateInfo
		templateName := templateInfo.template.Name

		runWithTimeout(t, testTimeout, templateName, parallel, func(t *testing.T) {
			if isBlackListedTest(templateName, cfg.skipped) {
				t.Skip("Skipping per BLACK_LISTED_TESTS")
			}

			t.Logf("Starting test run for %q", templateName)

			e := ptesting.NewEnvironment(t)
			t.Cleanup(func() { deleteIfNotFailed(e) })

			bench := guessBench(templateInfo.template)

			e.SetEnvVars(append(e.Env, bench.Env()...)...)

			pulumiNew(e, templateInfo.templatePath, bench.CommandArgs("pulumi-new")...)

			opts := base.
				With(detectOptionsFromTestingEnvironment(t, e)).
				With(integration.ProgramTestOptions{
					Dir:    e.RootPath,
					Config: cfg.config,
				}).
				With(bench.ProgramTestOptions())

			integration.ProgramTest(t, &opts)
		})
	}
}

func pulumiNew(e *ptesting.Environment, templatePath string, extraArgs ...string) {
	// Pulumi new expects a stack name or assumes dev, we generate
	// a random one here to prevent conflicts. Note that
	// ProgramTest will use its own stack, so we take care to
	// delete this one right away. There is a --generate-only
	// option but that opts out of installing dependencies, but we
	// want that to happen as part of pulumi new
	tempStack := (&integration.ProgramTestOptions{}).GetStackName().String()
	cmdArgs := append(
		[]string{"new", templatePath, "-f", "--yes", "-s", tempStack},
		extraArgs...,
	)
	e.RunCommand("pulumi", cmdArgs...)
	e.RunCommand("pulumi", "stack", "rm", tempStack, "--yes")
}

func detectOptionsFromTestingEnvironment(t *testing.T, e *ptesting.Environment) integration.ProgramTestOptions {
	path, err := workspace.DetectProjectPathFrom(e.RootPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	projinfo, err := workspace.LoadProject(path)
	assert.NoError(t, err)

	var prepareProject func(*engine.Projinfo) error
	switch rt := projinfo.Runtime.Name(); rt {
	case integration.NodeJSRuntime:
		// Default PrepareProject for Node uses yarn install
		// to install dependencies; template tests do not need
		// it because pulumi new already installs them with
		// npm, which is also what will happen on user
		// systems.
		prepareProject = func(*engine.Projinfo) error {
			return nil
		}
	default:
		prepareProject = nil // use default logic
	}
	return integration.ProgramTestOptions{
		PrepareProject: prepareProject,
	}
}

// deleteIfNotFailed deletes the files in the testing environment if the testcase has
// not failed. (Otherwise they are left to aid debugging.)
func deleteIfNotFailed(e *ptesting.Environment) {
	if _, found := os.LookupEnv("CI"); found {
		// Skip cleanup on CI, workaround for https://github.com/pulumi/pulumi/issues/9437
		return
	}
	if !e.T.Failed() {
		e.DeleteEnvironment()
	}
}

func isBlackListedTest(templateName string, backListedTests []string) bool {
	for _, blackListed := range backListedTests {
		if strings.Contains(templateName, blackListed) {
			return true
		}
	}

	return false
}
