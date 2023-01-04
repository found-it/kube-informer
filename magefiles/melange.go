//go:build mage
// +build mage

package main

import (
    "fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path/filepath"
    "os"
)

type Melange mg.Namespace

func (Melange) Keygen() error {
    client, err := GetClient()
    if err != nil {
        return err
    }

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

    args := []string{
        "run",
        "--rm",
	    "--privileged",
        "--volume",
        fmt.Sprintf("%s:/work", cwd),
        MELANGE_IMAGE,
        "keygen",
        MELANGE_PRIVATE_KEY,
    }

	return sh.RunV(client, args...)
}

func (Melange) Build() error {
    mg.Deps(Melange.Clean, Melange.Keygen)
    client, err := GetClient()
    if err != nil {
        return err
    }

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

    args := []string{
        "run",
        "--rm",
	    "--privileged",
        "--volume",
        fmt.Sprintf("%s:/work", cwd),
        MELANGE_IMAGE,
        "build",
        "melange.yaml",
        "--arch",
        "x86_64",
        "--repository-append",
        "packages",
        "--signing-key",
        MELANGE_PRIVATE_KEY,
    }

	return sh.RunV(client, args...)
}

// Cleans the packages directory
func (Melange) Clean() error {
    cwd, err := os.Getwd()
    if err != nil {
        return err
    }

    packages := filepath.Join(cwd, "packages")

    info, err := os.Stat(packages)

    if err == nil && info.IsDir() {
        fmt.Printf("Removing %s\n", packages)
        return os.RemoveAll(packages)
    } else if err != nil && !os.IsNotExist(err) {
        return err
    }

    return nil
}
