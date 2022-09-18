package testutils

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

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
			templatePath = path.Join(templateUrl, templateName)
		}

		infos = append(infos, TemplateInfo{
			Template:     t,
			TemplatePath: templatePath,
		})
	}
	return infos
}
