package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
	ptesting "github.com/pulumi/pulumi/sdk/v2/go/common/testing"
	"github.com/pulumi/pulumi/sdk/v2/go/common/workspace"
	"github.com/stretchr/testify/assert"
)

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

	base := integration.ProgramTestOptions{
		ExpectRefreshChanges:   true,
		Quick:                  true,
		SkipRefresh:            true,
		NoParallel:             true, // we mark tests as Parallel manually when instantiating
		UseAutomaticVirtualEnv: true,
	}

	// Retrieve the template repo.
	repo, err := workspace.RetrieveTemplates("", false /*offline*/, workspace.TemplateKindPulumiProject)
	assert.NoError(t, err)
	defer assert.NoError(t, repo.Delete())

	// List the templates from the repo.
	templates, err := repo.Templates()
	assert.NoError(t, err)

	blackListed := strings.Split(blackListedTests, ",")

	for _, template := range templates {
		templateName := template.Name
		t.Run(templateName, func(t *testing.T) {
			t.Parallel()
			if isBlackListedTest(templateName, blackListed) {
				t.Skipf("Skipping template test %s", templateName)
				return
			}

			t.Logf("Starting test run for %q", templateName)

			e := ptesting.NewEnvironment(t)
			defer deleteIfNotFailed(e)

			e.RunCommand("pulumi", "new", templateName, "-f", "--yes", "-s", "template-test")

			path, err := workspace.DetectProjectPathFrom(e.RootPath)
			assert.NoError(t, err)
			assert.NotEmpty(t, path)

			_, err = workspace.LoadProject(path)
			assert.NoError(t, err)

			example := base.With(integration.ProgramTestOptions{
				Dir: e.RootPath,
				Config: map[string]string{
					"aws:region":        awsRegion,
					"azure:environment": azureEnviron,
					"azure:location":    azureLocation,
					"gcp:project":       gcpProject,
					"gcp:region":        gcpRegion,
					"gcp:zone":          gcpZone,
					"cloud:provider":    "aws",
				},
			})

			integration.ProgramTest(t, &example)
		})
	}
}

// deleteIfNotFailed deletes the files in the testing environment if the testcase has
// not failed. (Otherwise they are left to aid debugging.)
func deleteIfNotFailed(e *ptesting.Environment) {
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
