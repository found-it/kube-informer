//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
)

type Inform mg.Namespace

func getBinPath() (string, error) {
	var (
		bin     string = "inform"
		binDir  string = "bin"
		binPath string = filepath.Join(binDir, bin)
	)
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, binPath), nil
}

// Builds the inform binary in ./bin
func (Inform) Build() error {
	mg.Deps(Inform.Clean)
	binPath, err := getBinPath()
	if err != nil {
		return err
	}

    fmt.Printf("[inform:build] Creating %s\n", filepath.Dir(binPath))
	err = os.Mkdir(filepath.Dir(binPath), 0744)
	if err != nil && !os.IsExist(err) {
		return err
	}

	args := []string{
		"build",
		"-a",
		"-installsuffix",
		"cgo",
		"-o",
		binPath,
		".",
	}

    env := map[string]string{
        "CGO_ENABLED": "0",
        "GOOS": "linux",
        "GOARCH": "amd64",
    }

    fmt.Printf("[inform:build] Building %s\n", binPath)
	return sh.RunWith(env, "go", args...)
}

// Builds and runs the inform binary out of ./bin
func (Inform) Run() error {
    mg.Deps(Inform.Build)
	binPath, err := getBinPath()
	if err != nil {
		return err
	}

    fmt.Printf("[inform:clean] Running %s\n", binPath)
	return sh.RunV(binPath)
}

// Cleans the inform binary out of ./bin
func (Inform) Clean() error {
	binPath, err := getBinPath()
	if err != nil {
		return err
	}

    fmt.Printf("[inform:clean] Cleaning %s\n", binPath)
	return os.RemoveAll(filepath.Dir(binPath))
}
