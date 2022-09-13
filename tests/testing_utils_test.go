package tests

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type templateInfo struct {
	template     workspace.Template
	templatePath string
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
