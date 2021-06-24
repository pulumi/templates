// TODO a lot of duplication here against pulumi/examples
// performance_test.go, simplify somehow.

package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi-trace-tool/traces"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type bench struct {
	name     string
	provider string
	runtime  string
	language string
}

func TestMain(m *testing.M) {
	code := m.Run()

	dir := tracingDir()
	if dir != "" {
		// After all tests run with tracing, compute metrics
		// on the entire set.
		err := computeMetrics(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(code)
}

func tracingDir() string {
	return os.Getenv("PULUMI_TRACING_DIR")
}

func tracingOpts(t *testing.T, benchmark bench) integration.ProgramTestOptions {
	dir := tracingDir()

	if dir != "" {
		return integration.ProgramTestOptions{
			Env: []string{
				"PULUMI_TRACING_TAG_REPO=pulumi/templates",
				fmt.Sprintf("PULUMI_TRACING_TAG_BENCHMARK_NAME=%s", benchmark.name),
				fmt.Sprintf("PULUMI_TRACING_TAG_BENCHMARK_PROVIDER=%s", benchmark.provider),
				fmt.Sprintf("PULUMI_TRACING_TAG_BENCHMARK_RUNTIME=%s", benchmark.runtime),
				fmt.Sprintf("PULUMI_TRACING_TAG_BENCHMARK_LANGUAGE=%s", benchmark.language),
				"PULUMI_TRACING_MEMSTATS_POLL_INTERVAL=100ms",
			},
			Tracing: fmt.Sprintf("file:%s",
				filepath.Join(dir, fmt.Sprintf("%s-{command}.trace", benchmark.name))),
		}
	}

	return integration.ProgramTestOptions{}
}

func computeMetrics(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	defer os.Chdir(cwd)

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}

	var traceFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".trace") {
			traceFiles = append(traceFiles, f.Name())
		}
	}

	if err := traces.ToCsv(traceFiles, "traces.csv", "filename"); err != nil {
		return err
	}

	f, err := os.Create("metrics.csv")
	if err != nil {
		return err
	}

	if err := traces.Metrics("traces.csv", "filename", f); err != nil {
		return err
	}

	return nil
}

// additional helpers specific to this repo (templates)

func guessBench(template workspace.Template) bench {
	name := filepath.Base(template.Dir)
	return bench{
		name:     name,
		provider: guessProvider(template),
		language: guessLanguage(template),
		runtime:  guessRuntime(template),
	}
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
