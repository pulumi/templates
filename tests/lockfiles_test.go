package tests

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var otherPackageManagerLockfiles = []string{"yarn.lock", "pnpm-lock.yaml", "bun.lockb"}

func TestLockfileMatchesDeclaredRuntime(t *testing.T) {
	root := repoRoot(t)

	for _, project := range nodeProjectDirs(t, root) {
		t.Run(project, func(t *testing.T) {
			carries, forbids := "package-lock.json", "bun.lock"
			if declaredRuntime(t, filepath.Join(root, project)) == "bun" {
				carries, forbids = "bun.lock", "package-lock.json"
			}

			assert.FileExists(t, filepath.Join(root, project, carries),
				"%s declares a runtime that installs with the package manager behind %s, so that is the lockfile it has to carry",
				project, carries)
			assert.NoFileExists(t, filepath.Join(root, project, forbids),
				"the Pulumi CLI picks a project's package manager by which lockfile it finds, so %s in %s hands the template to a package manager it never declared",
				forbids, project)
		})
	}
}

func TestLockfileHasAPackageJSONBesideIt(t *testing.T) {
	root := repoRoot(t)

	for _, lockfile := range committedLockfilePaths(t, root) {
		assert.FileExists(t, filepath.Join(root, filepath.Dir(lockfile), "package.json"),
			"%s has no package.json beside it, so nothing regenerates it and it will drift", lockfile)
	}
}

func TestLockfileIsNeverFromAnotherPackageManager(t *testing.T) {
	root := repoRoot(t)

	for _, project := range nodeProjectDirs(t, root) {
		for _, name := range otherPackageManagerLockfiles {
			assert.NoFileExists(t, filepath.Join(root, project, name),
				"%s belongs to a package manager no template declares, and its presence would hand %s to that package manager",
				name, project)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs("..")
	require.NoError(t, err)
	return root
}

func declaredRuntime(t *testing.T, dir string) string {
	t.Helper()

	contents, err := os.ReadFile(filepath.Join(dir, "Pulumi.yaml"))
	if os.IsNotExist(err) {
		return ""
	}
	require.NoError(t, err)

	var project struct {
		Runtime yaml.Node `yaml:"runtime"`
	}
	require.NoError(t, yaml.Unmarshal(contents, &project))

	switch project.Runtime.Kind {
	case yaml.ScalarNode:
		return project.Runtime.Value
	case yaml.MappingNode:
		var runtime struct {
			Name string `yaml:"name"`
		}
		require.NoError(t, project.Runtime.Decode(&runtime))
		return runtime.Name
	default:
		return ""
	}
}

func nodeProjectDirs(t *testing.T, root string) []string {
	t.Helper()

	var projects []string
	require.NoError(t, filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && nonTemplateDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() != "package.json" {
			return nil
		}

		dir := filepath.Dir(path)
		if dir == root {
			return nil
		}
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			return err
		}
		projects = append(projects, filepath.ToSlash(rel))
		return nil
	}))

	require.NotEmpty(t, projects, "found no Node projects, so this test is not testing anything")
	sort.Strings(projects)
	return projects
}

func committedLockfilePaths(t *testing.T, root string) []string {
	t.Helper()

	var lockfiles []string
	require.NoError(t, filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
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
		lockfiles = append(lockfiles, filepath.ToSlash(rel))
		return nil
	}))

	require.NotEmpty(t, lockfiles, "found no lockfiles, so this test is not testing anything")
	sort.Strings(lockfiles)
	return lockfiles
}
