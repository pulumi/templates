package tests

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type templateInfo struct {
	template     workspace.Template
	templatePath string
}

type templateTestConfig struct {
	config      map[string]string
	templateUrl string
	skipped     []string
}

func newTemplateTestConfigFromEnv() templateTestConfig {
	skipped := []string{}
	if l := os.Getenv("BLACK_LISTED_TESTS"); l != "" {
		skipped = strings.Split(l, ",")
	}

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
	if loc := os.Getenv("PULUMI_TEMPLATE_LOCATION"); loc != "" {
		templateUrl = loc
	}

	return templateTestConfig{
		skipped:     skipped,
		templateUrl: templateUrl,
		config: map[string]string{
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
	}
}

func findAllTemplates(t *testing.T, templateUrl string) []templateInfo {
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

	infos := []templateInfo{}
	for _, t := range templates {
		templateName := t.Name
		templatePath := templateName
		if templateUrl != "" {
			templatePath = path.Join(templateUrl, templateName)
		}

		infos = append(infos, templateInfo{
			template:     t,
			templatePath: templatePath,
		})
	}
	return infos
}

func runWithTimeout(
	t *testing.T,
	timeout time.Duration,
	name string,
	prepare func(*testing.T),
	run func(*testing.T),
) {
	t.Run(name, func(t *testing.T) {
		prepare(t)
		timeoutEvent := time.After(timeout)
		done := make(chan bool)
		go func() {
			defer func() {
				done <- true
			}()
			run(t)
		}()
		select {
		case <-timeoutEvent:
			t.Fatalf("%s timed out after %s", name, timeout)
		case <-done:
			return
		}
	})
}

func parallel(t *testing.T) {
	t.Parallel()
}
