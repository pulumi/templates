package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type TemplateInfo struct {
	Template     workspace.Template
	TemplatePath string
}

func FindAllTemplates(t *testing.T, templateUrl string) []TemplateInfo {
	// Retrieve the template repo.
	repo, err := workspace.RetrieveTemplates(t.Context(), templateUrl, false /*offline*/, workspace.TemplateKindPulumiProject)
	assert.NoError(t, err)
	t.Cleanup(func() {
		err := repo.Delete()
		assert.NoError(t, err, "Error cleaning up repository after deletion.")
	})

	// List the templates from the repo.
	templates, err := repo.Templates()
	assert.NoError(t, err)

	infos := []TemplateInfo{}
	for _, t := range templates {
		templateName := t.Name
		templatePath := templateName
		if templateUrl != "" {
			templatePath = filepath.Join(templateUrl, templateName)
		}

		infos = append(infos, TemplateInfo{
			Template:     t,
			TemplatePath: templatePath,
		})
	}
	return infos
}

// UpdateOptions returns the set of integration.ProgramTestOptions that should be applied for the
// given template.
func UpdateOptions(t testing.TB, templateInfo TemplateInfo) integration.ProgramTestOptions {

	// For templates marked important, we test the full end to end experience to ensure
	// updates succeed and subsequent operations produce no changes.
	templatePath := filepath.Join(templateInfo.TemplatePath, "Pulumi.yaml")
	templateBytes, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	var template map[string]any
	require.NoError(t, yaml.Unmarshal(templateBytes, &template))

	skipFullTest := true
	if important, ok := template["template"].(map[string]any)["important"]; ok {
		skipFullTest = !important.(bool)
	}

	return integration.ProgramTestOptions{
		// Skip running a full update.
		SkipUpdate: skipFullTest,
		// Skip running a refresh after the update, expecting no changes.
		SkipRefresh: skipFullTest,
		// Skip running a preview after the update, expecting no changes.
		SkipPreview: skipFullTest,
		// Skip running a stack export and import after the update.
		SkipExportImport: skipFullTest,
	}
}
