package tests

import (
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/templates/tests/internal/testutils"
)

func TestTemplates(t *testing.T) {
	cfg := testutils.NewTemplateTestConfigFromEnv(t, testutils.SKIPPED_TESTS)

	for _, templateInfo := range testutils.FindAllTemplates(t, cfg.TemplateUrl) {
		templateInfo := templateInfo
		templateName := templateInfo.Template.Name

		t.Run(templateName, func(t *testing.T) {
			cfg.PossiblySkip(t, templateInfo)
			t.Parallel()

			e := testutils.NewEnvironment(t, cfg)
			testutils.PulumiNew(e, templateInfo.TemplatePath)

			opts := integration.ProgramTestOptions{
				Dir:                    e.RootPath,
				Config:                 cfg.Config,
				DestroyOnCleanup:       true,
				UseAutomaticVirtualEnv: true,
				NoParallel:             true, // Called before
				PrepareProject:         testutils.PrepareProject(t, e),
				RequireService:         true,
				InstallDevReleases:     true,
			}.With(testutils.UpdateOptions(t, templateInfo))

			integration.ProgramTest(t, &opts)
		})
	}
}
