package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	root := ".."
	for i, f := range findPackageJsonFiles() {
		pkgJson := filepath.Join(root, "tests", "npm-lockfiles", f)
		if err := rewriteName(fmt.Sprintf("template%d", i),
			filepath.Join(root, f),
			pkgJson); err != nil {
			log.Fatal(err)
		}
		if err := computeLock(root, pkgJson); err != nil {
			log.Fatal(err)
		}
		if err := os.Remove(pkgJson); err != nil {
			log.Fatal(err)
		}
	}
}

func computeLock(root, pkgJson string) error {
	d := filepath.Dir(pkgJson)
	fmt.Printf("npm install --package-lock-only %v\n", d)
	cmd := exec.Command("npm", "install", "--package-lock-only")
	cmd.Dir = d
	return cmd.Run()
}

func rewriteName(name, src, dest string) error {
	srcBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var srcData map[string]interface{}
	if err := json.Unmarshal(srcBytes, &srcData); err != nil {
		return err
	}
	srcData["name"] = name
	destBytes, err := json.Marshal(srcData)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(dest, destBytes, 0600)
}

func findPackageJsonFiles() []string {
	var buf bytes.Buffer
	cmd := exec.Command("git", "ls-files")
	cmd.Stdout = &buf
	cmd.Dir = ".."

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	files := strings.Split(buf.String(), "\n")
	var out []string
	for _, f := range files {
		if strings.HasSuffix(f, "package.json") &&
			!strings.Contains(f, "npm-lockfiles") {
			out = append(out, f)
		}
	}

	sort.Strings(out)
	return out
}
