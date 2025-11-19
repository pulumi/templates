package testutils

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const SKIPPED_TESTS = "SKIPPED_TESTS"
const SKIPPED_BENCHMARKS = "SKIPPED_BENCHMARKS"

type TemplateTestConfig struct {
	Config      map[string]string
	TemplateUrl string
	Skipped     []string
	SkipEnvVar  string
}

func NewTemplateTestConfigFromEnv(t testing.TB, skipEnvVar string) TemplateTestConfig {
	skipped := []string{}
	if l := os.Getenv(skipEnvVar); l != "" {
		skipped = strings.Split(l, ",")
		t.Logf("Reading test skip pattern from env var %s", skipEnvVar)
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-west-1"
		t.Log("Defaulting AWS_REGION to 'us-west-1'.  You can override using the AWS_REGION environment variable")
	}
	azureEnviron := os.Getenv("ARM_ENVIRONMENT")
	if azureEnviron == "" {
		azureEnviron = "public"
		t.Log("Defaulting ARM_ENVIRONMENT to 'public'.  You can override using the ARM_ENVIRONMENT variable")
	}
	azureLocation := os.Getenv("ARM_LOCATION")
	if azureLocation == "" {
		azureLocation = "westus"
		t.Log("Defaulting ARM_LOCATION to 'westus'.  You can override using the ARM_LOCATION variable")
	}
	gcpProject := os.Getenv("GOOGLE_PROJECT")
	if gcpProject == "" {
		gcpProject = "pulumi-ci-gcp-provider"
		t.Log("Defaulting GOOGLE_PROJECT to 'pulumi-ci-gcp-provider'." +
			"You can override using the GOOGLE_PROJECT variable")
	}
	gcpRegion := os.Getenv("GOOGLE_REGION")
	if gcpRegion == "" {
		gcpRegion = "us-central1"
		t.Log("Defaulting GOOGLE_REGION to 'us-central1'.  You can override using the GOOGLE_REGION variable")
	}
	gcpZone := os.Getenv("GOOGLE_ZONE")
	if gcpZone == "" {
		gcpZone = "us-central1-a"
		t.Log("Defaulting GOOGLE_ZONE to 'us-central1-a'.  You can override using the GOOGLE_ZONE variable")
	}

	// by default, we want to test the normal template url path
	// if we have a specific template location set then we should
	// use that in our tests
	templateUrl := ""
	if loc := os.Getenv("PULUMI_TEMPLATE_LOCATION"); loc != "" {
		templateUrl = loc
		t.Log("Using templates from PULUMI_TEMPLATE_LOCATION=%s\n", loc)
	}

	return TemplateTestConfig{
		SkipEnvVar:  skipEnvVar,
		Skipped:     skipped,
		TemplateUrl: templateUrl,
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
	}
}

func (cfg TemplateTestConfig) IsSkipped(info TemplateInfo) bool {
	for _, s := range cfg.Skipped {
		if strings.Contains(info.Template.Name, s) {
			return true
		}
	}
	return false
}

func (cfg TemplateTestConfig) PossiblySkip(t *testing.T, info TemplateInfo) {
	if cfg.IsSkipped(info) {
		t.Skip(fmt.Sprintf("Skipping per %s", cfg.SkipEnvVar))
	}
}
