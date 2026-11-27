package tests

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/backend"
	cmdBackend "github.com/pulumi/pulumi/pkg/v3/cmd/pulumi/backend"
	"github.com/pulumi/pulumi/pkg/v3/cmd/pulumi/templatecmd"
	pkgWorkspace "github.com/pulumi/pulumi/pkg/v3/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag/colors"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var committedLockfileNames = map[string]bool{
	"package-lock.json": true,
	"bun.lock":          true,
}

var nonTemplateDirs = map[string]bool{
	"generator":    true,
	"metadata":     true,
	"scripts":      true,
	"tests":        true,
	"node_modules": true,
	".git":         true,
	".jj":          true,
}

func TestTemplatePublishIncludesLockfiles(t *testing.T) {
	root, err := filepath.Abs("..")
	require.NoError(t, err)

	lockfiles := lockfilesByTemplate(t, root)
	require.NotEmpty(t, lockfiles, "found no lockfiles, so this test is not testing anything")

	templates := make([]string, 0, len(lockfiles))
	for template := range lockfiles {
		templates = append(templates, template)
	}
	sort.Strings(templates)

	for _, template := range templates {
		t.Run(template, func(t *testing.T) {
			published := publishTemplate(t, filepath.Join(root, template), template)

			for _, lockfile := range lockfiles[template] {
				assert.Contains(t, published, lockfile,
					"pulumi template publish drops %s/%s, most likely because a .gitignore entry excludes it, so it never reaches the users who install this template",
					template, lockfile)
			}
		})
	}
}

func publishTemplate(t *testing.T, templateDir, name string) []string {
	t.Helper()

	var archived []byte
	registry := &backend.MockCloudRegistry{
		PublishTemplateF: func(_ context.Context, op apitype.TemplatePublishOp) error {
			var err error
			archived, err = io.ReadAll(op.Archive)
			return err
		},
	}
	fakeBackend := &backend.MockBackend{
		GetDefaultOrgF:    func(context.Context) (string, error) { return "pulumi", nil },
		GetCloudRegistryF: func() (backend.CloudRegistry, error) { return registry, nil },
	}

	restore := cmdBackend.DefaultLoginManager
	cmdBackend.DefaultLoginManager = &fakeLoginManager{backend: fakeBackend}
	t.Cleanup(func() { cmdBackend.DefaultLoginManager = restore })

	cmd := templatecmd.NewTemplateCmd()
	cmd.SetArgs([]string{"publish", templateDir, "--name", name, "--version", "0.0.0"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	require.NoError(t, cmd.ExecuteContext(t.Context()))
	require.NotEmpty(t, archived, "the publish command uploaded no archive")

	return pathsInArchive(t, archived)
}

type fakeLoginManager struct {
	backend backend.Backend
}

func (m *fakeLoginManager) Current(
	context.Context, pkgWorkspace.Context, diag.Sink, string, *workspace.Project, bool,
) (backend.Backend, error) {
	return m.backend, nil
}

func (m *fakeLoginManager) Login(
	context.Context, pkgWorkspace.Context, diag.Sink, string, *workspace.Project, bool, bool, colors.Colorization,
) (backend.Backend, error) {
	return m.backend, nil
}

func (m *fakeLoginManager) LoginFromAuthContext(
	context.Context, diag.Sink, string, *workspace.Project, bool, bool, workspace.AuthContext,
) (backend.Backend, error) {
	return m.backend, nil
}

func lockfilesByTemplate(t *testing.T, root string) map[string][]string {
	t.Helper()

	lockfiles := map[string][]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && nonTemplateDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !committedLockfileNames[entry.Name()] {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		segments := strings.Split(filepath.ToSlash(rel), "/")
		if len(segments) < 2 {
			return nil
		}
		template := segments[0]
		lockfiles[template] = append(lockfiles[template], strings.Join(segments[1:], "/"))
		return nil
	})
	require.NoError(t, err)

	return lockfiles
}

func pathsInArchive(t *testing.T, tgz []byte) []string {
	t.Helper()

	gz, err := gzip.NewReader(bytes.NewReader(tgz))
	require.NoError(t, err)
	defer gz.Close()

	var names []string
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		names = append(names, strings.TrimPrefix(header.Name, "./"))
	}
	return names
}
