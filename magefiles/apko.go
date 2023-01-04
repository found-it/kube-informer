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

type APKO mg.Namespace

func (APKO) Build() error {
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
        "--workdir",
        "/work",
        APKO_IMAGE,
        "build",
        "--debug",
        "image.yaml",
        "jpetersenames/inform:test",
        "inform.tar",
        "--keyring-append",
        MELANGE_PUBLIC_KEY,
    }

	return sh.RunV(client, args...)
}
// Cleans the apko artifacts
func (APKO) Clean() error {
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
