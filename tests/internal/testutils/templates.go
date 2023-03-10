package testutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type TemplateInfo struct {
	Template     workspace.Template
	TemplatePath string
}

func FindAllTemplates(t *testing.T, templateUrl string) []TemplateInfo {
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
func UpdateOptions(templateName string) integration.ProgramTestOptions {

	// For our core Getting Started templates, we test the full end to end experience,
	// ensuring updates succeed and subsequent operations produce no changes.
	skipFullTest := !isGettingStartedTemplate(templateName)

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

// isGettingStartedTemplate returns whether the specified template is one that is used in
// our Getting Started flows.
func isGettingStartedTemplate(templateName string) bool {
	coreStarterTemplateNames := []string{
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

	for _, s := range coreStarterTemplateNames {
		if s == templateName {
			return true
		}
	}
	return false
}
