package perf

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi-trace-tool/traces"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"

	"github.com/pulumi/templates/v2/internal/testutils"
)

const testTimeout = 60 * time.Minute

func TestTemplatePerf(t *testing.T) {
	cfg := testutils.NewTemplateTestConfigFromEnv()

	for _, templateInfo := range testutils.FindAllTemplates(t, cfg.TemplateUrl) {
		templateInfo := templateInfo
		templateName := templateInfo.Template.Name

		prepare := func(t *testing.T) {
			cfg.PossiblySkip(t, templateInfo)
		}

		testutils.RunWithTimeout(t, testTimeout, templateName, prepare, func(t *testing.T) {
			t.Logf("Starting test run for %q", templateName)

			bench := guessBench(templateInfo.Template)

			e := testutils.NewEnvironment(t, cfg)
			e.SetEnvVars(append(e.Env, bench.Env()...)...)

			testutils.PulumiNew(e, templateInfo.TemplatePath,
				bench.CommandArgs("pulumi-new")...)

			opts := integration.ProgramTestOptions{
				Dir:                    e.RootPath,
				Config:                 cfg.Config,
				ExpectRefreshChanges:   true,
				Quick:                  false, // true skews measurements
				SkipRefresh:            true,
				NoParallel:             true, // minimize interference
				DestroyOnCleanup:       true,
				UseAutomaticVirtualEnv: true,
				PrepareProject:         testutils.PrepareProject(t, e),
			}.With(bench.ProgramTestOptions())

			integration.ProgramTest(t, &opts)
		})
	}

	if err := traces.ComputeMetrics(); err != nil {
		t.Fatalf("ComputeMetrics failed: %v", err)
		t.FailNow()
	}

}

func guessBench(template workspace.Template) traces.Benchmark {
	b := traces.NewBenchmark(template.Name)
	b.Provider = guessProvider(template)
	b.Language = guessLanguage(template)
	b.Runtime = guessRuntime(template)
	b.Repository = "pulumi/templates"
	return b
}

func guessLanguage(template workspace.Template) string {
	languages := map[string]string{
		"fsharp":      "fsharp",
		"csharp":      "csharp",
		"py":          "python",
		"go":          "go",
		"javascript":  "javascript",
		"typescript":  "typescript",
		"visualbasic": "visualbasic",
	}
	for _, lang := range languages {
		if strings.Contains(template.Name, lang) {
			return lang
		}
	}
	return ""
}

func guessRuntime(template workspace.Template) string {
	proj, err := workspace.LoadProject(filepath.Join(template.Dir, "Pulumi.yaml"))
	if err != nil {
		return ""
	}

	return proj.Runtime.Name()
}

func guessProvider(template workspace.Template) string {
	if strings.Contains(template.Name, "azure-classic") {
		return "azure-classic"
	}

	if strings.Contains(template.Name, "azure") {
		return "azure"
	}

	if strings.Contains(template.Name, "google-native") {
		return "google-native"
	}

	if strings.Contains(template.Name, "gcp") {
		return "gcp"
	}

	if strings.Contains(template.Name, "aws") {
		return "aws"
	}

	if strings.Contains(template.Name, "kubernetes") {
		return "kubernetes"
	}

	if strings.Contains(template.Name, "openstack") {
		return "openstack"
	}

	if strings.Contains(template.Name, "linode") {
		return "linode"
	}

	if strings.Contains(template.Name, "digitalocean") {
		return "digitalocean"
	}

	if strings.Contains(template.Name, "alicloud") {
		return "alicloud"
	}

	if strings.Contains(template.Name, "equinix-metal") {
		return "equinix-metal"
	}

	return ""
}
