package perf

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-trace-tool/traces"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"

	"github.com/pulumi/templates/tests/internal/testutils"
)

func TestTemplatePerf(t *testing.T) {
	if !traces.IsTracingEnabled() {
		t.Fatalf("Required environment variable not set: %s", traces.TRACING_DIR_ENV_VAR)
	}

	cfg := testutils.NewTemplateTestConfigFromEnv(t, testutils.SKIPPED_BENCHMARKS)

	for _, templateInfo := range testutils.FindAllTemplates(t, cfg.TemplateUrl) {
		templateInfo := templateInfo

		t.Run(templateInfo.Template.Name, func(t *testing.T) {
			cfg.PossiblySkip(t, templateInfo)
			t.Parallel()

			t.Run("prewarm", func(t *testing.T) {
				prewarm(t, cfg, templateInfo)
			})

			t.Run("benchmark", func(t *testing.T) {
				benchmark(t, cfg, templateInfo)
			})
		})
	}

	if err := traces.ComputeMetrics(); err != nil {
		t.Fatalf("traces.ComputeMetrics() failed: %v", err)
	}
}

// Prewarming runs preview only to make sure all needed plugins are
// downloaded so that these downloads do not skew measurements.
func prewarm(t *testing.T, cfg testutils.TemplateTestConfig, templateInfo testutils.TemplateInfo) {
	e := testutils.NewEnvironment(t, cfg)
	testutils.PulumiNew(e, templateInfo.TemplatePath)

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:                    e.RootPath,
		Config:                 cfg.Config,
		NoParallel:             true,
		PrepareProject:         testutils.PrepareProject(t, e),
		SkipRefresh:            true,
		SkipEmptyPreviewUpdate: true,
		SkipExportImport:       true,
		SkipUpdate:             true,
	})
}

func benchmark(t *testing.T, cfg testutils.TemplateTestConfig, templateInfo testutils.TemplateInfo) {
	templateName := templateInfo.Template.Name

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
		SkipExportImport:       true, // save time on CI
		NoParallel:             true, // minimize interference
		DestroyOnCleanup:       true,
		UseAutomaticVirtualEnv: true,
		PrepareProject:         testutils.PrepareProject(t, e),
		RequireService:         true,
	}.With(bench.ProgramTestOptions())

	integration.ProgramTest(t, &opts)
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
