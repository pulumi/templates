package tests

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-trace-tool/traces"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

func TestMain(m *testing.M) {
	code := m.Run()

	// If tracing is enabled, compute metrics after running all the tests.
	err := traces.ComputeMetrics()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func isTracingEnabled() bool {
	return traces.IsTracingEnabled()
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
