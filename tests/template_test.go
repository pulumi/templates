package tests

import (
	"testing"
	"time"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"

	"github.com/pulumi/templates/v2/internal/testutils"
)

const testTimeout = 60 * time.Minute

func TestTemplates(t *testing.T) {
	cfg := testutils.NewTemplateTestConfigFromEnv(testutils.SKIPPED_TESTS)

	for _, templateInfo := range testutils.FindAllTemplates(t, cfg.TemplateUrl) {
		templateInfo := templateInfo
		templateName := templateInfo.Template.Name

		prepare := func(t *testing.T) {
			cfg.PossiblySkip(t, templateInfo)
			t.Parallel()
		}

		templatesToTest := []string{
			"aws-typescript",
			"aws-javascript",
			"aws-python",
			"aws-go",
			"aws-csharp",
			"aws-java",
			"aws-yaml",

			"azure-typescript",
			"azure-javascript",
			"azure-python",
			"azure-go",
			"azure-csharp",
			"azure-java",
			"azure-yaml",

			"gcp-typescript",
			"gcp-javascript",
			"gcp-python",
			"gcp-go",
			"gcp-csharp",
			"gcp-java",
			"gcp-yaml",

			"kubernetes-typescript",
			"kubernetes-javascript",
			"kubernetes-python",
			"kubernetes-go",
			"kubernetes-csharp",
			"kubernetes-java",
			"kubernetes-yaml",
		}

		testutils.RunWithTimeout(t, testTimeout, templateName, prepare, func(t *testing.T) {
			t.Logf("Starting test run for %q", templateName)

			e := testutils.NewEnvironment(t, cfg)
			testutils.PulumiNew(e, templateInfo.TemplatePath)

			isStarterTemplate := !testutils.ListContains(templatesToTest, templateName)

			integration.ProgramTest(t, &integration.ProgramTestOptions{
				Dir:    e.RootPath,
				Config: cfg.Config,

				NoParallel:             true, // marked Parallel by prepare
				DestroyOnCleanup:       true,
				UseAutomaticVirtualEnv: true,
				PrepareProject:         testutils.PrepareProject(t, e),
				RequireService:         true,

				// Always run full updates for core starter templates.
				SkipUpdate: !testutils.ListContains(templatesToTest, templateName),

				// Ensure there's no diff on subsequent updates.
				SkipRefresh:      !isStarterTemplate,
				SkipPreview:      !isStarterTemplate,
				SkipExportImport: !isStarterTemplate,
			})
		})
	}
}
