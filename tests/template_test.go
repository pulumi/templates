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

		testutils.RunWithTimeout(t, testTimeout, templateName, prepare, func(t *testing.T) {
			t.Logf("Starting test run for %q", templateName)

			e := testutils.NewEnvironment(t, cfg)
			testutils.PulumiNew(e, templateInfo.TemplatePath)

			opts := integration.ProgramTestOptions{
				Dir:                    e.RootPath,
				Config:                 cfg.Config,
				NoParallel:             true, // marked Parallel by prepare
				DestroyOnCleanup:       true,
				UseAutomaticVirtualEnv: true,
				PrepareProject:         testutils.PrepareProject(t, e),
				RequireService:         true,
				InstallDevReleases:     true,
			}.With(testutils.UpdateOptions(templateInfo))

			integration.ProgramTest(t, &opts)
		})
	}
}
