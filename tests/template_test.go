package tests

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	ptesting "github.com/pulumi/pulumi/sdk/v3/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 1 * time.Minute

func TestTemplates(t *testing.T) {
	blackListedTests := os.Getenv("BLACK_LISTED_TESTS")

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-west-1"
		fmt.Println("Defaulting AWS_REGION to 'us-west-1'.  You can override using the AWS_REGION environment variable")
	}
	azureEnviron := os.Getenv("ARM_ENVIRONMENT")
	if azureEnviron == "" {
		azureEnviron = "public"
		fmt.Println("Defaulting ARM_ENVIRONMENT to 'public'.  You can override using the ARM_ENVIRONMENT variable")
	}
	azureLocation := os.Getenv("ARM_LOCATION")
	if azureLocation == "" {
		azureLocation = "westus"
		fmt.Println("Defaulting ARM_LOCATION to 'westus'.  You can override using the ARM_LOCATION variable")
	}
	gcpProject := os.Getenv("GOOGLE_PROJECT")
	if gcpProject == "" {
		gcpProject = "pulumi-ci-gcp-provider"
		fmt.Println("Defaulting GOOGLE_PROJECT to 'pulumi-ci-gcp-provider'." +
			"You can override using the GOOGLE_PROJECT variable")
	}
	gcpRegion := os.Getenv("GOOGLE_REGION")
	if gcpRegion == "" {
		gcpRegion = "us-central1"
		fmt.Println("Defaulting GOOGLE_REGION to 'us-central1'.  You can override using the GOOGLE_REGION variable")
	}
	gcpZone := os.Getenv("GOOGLE_ZONE")
	if gcpZone == "" {
		gcpZone = "us-central1-a"
		fmt.Println("Defaulting GOOGLE_ZONE to 'us-central1-a'.  You can override using the GOOGLE_ZONE variable")
	}

	// by default, we want to test the normal template url path
	// if we have a specific template location set then we should
	// use that in our tests
	templateUrl := ""
	specificTemplate := os.Getenv("PULUMI_TEMPLATE_LOCATION")
	if specificTemplate != "" {
		templateUrl = specificTemplate
	}

	base := integration.ProgramTestOptions{
		ExpectRefreshChanges:   true,
		Quick:                  true,
		SkipRefresh:            true,
		NoParallel:             true, // we mark tests as Parallel manually when instantiating
		DestroyOnCleanup:       true,
		UseAutomaticVirtualEnv: true,
	}

	// Retrieve the template repo.
	repo, err := workspace.RetrieveTemplates(templateUrl, false /*offline*/, workspace.TemplateKindPulumiProject)
	assert.NoError(t, err)
	t.Cleanup(func() {
		err := repo.Delete()
		assert.NoError(t, err, "Error cleaning up repository after deletion.")
	})

	// List the templates from the repo.
	templates, err := repo.Templates()
	assert.NoError(t, err)

	blackListed := strings.Split(blackListedTests, ",")

	for _, template := range templates {
		template := template
		templateName := template.Name

		if isBlackListedTest(templateName, blackListed) {
			t.Logf("Skipping template test %s", templateName)
			continue
		}

		e := ptesting.NewEnvironment(t)
		t.Cleanup(func() { deleteIfNotFailed(e) })

		runWithTimeout(t, testTimeout, templateName, func(t *testing.T) {
			t.Parallel()

			t.Logf("Starting test run for %q", templateName)

			bench := guessBench(template)

			e.SetEnvVars(append(e.Env, bench.Env()...))

			templatePath := templateName
			if templateUrl != "" {
				templatePath = path.Join(templateUrl, templateName)
			}

			cmdArgs := append(
				[]string{"new", templatePath, "-f", "--yes", "-s", "template-test"},
				bench.CommandArgs("pulumi-new")...,
			)

			e.RunCommand("pulumi", cmdArgs...)

			path, err := workspace.DetectProjectPathFrom(e.RootPath)
			assert.NoError(t, err)
			assert.NotEmpty(t, path)

			_, err = workspace.LoadProject(path)
			assert.NoError(t, err)

			example := base.With(integration.ProgramTestOptions{
				Dir: e.RootPath,
				Config: map[string]string{
					"aws:region":            awsRegion,
					"azure:environment":     azureEnviron,
					"azure:location":        azureLocation,
					"azure-native:location": azureLocation,
					"gcp:project":           gcpProject,
					"gcp:region":            gcpRegion,
					"gcp:zone":              gcpZone,
					"google-native:project": gcpProject,
					"google-native:region":  gcpRegion,
					"google-native:zone":    gcpZone,
					"cloud:provider":        "aws",
				},
			}).With(bench.ProgramTestOptions())

			integration.ProgramTest(t, &example)
		})
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
