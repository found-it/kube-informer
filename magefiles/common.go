package main

import (
	"fmt"
	"os/exec"
)

const WOLFI_PACKAGES_URL = "https://packages.wolfi.dev/os"
const WOLFI_PACKAGES_KEY_URL = "https://packages.wolfi.dev/os/wolfi-signing.rsa.pub"

const MELANGE_IMAGE = "cgr.dev/chainguard/melange@sha256:7dfc7861d0946c04b23bfd7bdd90110b3cb40711fe760f0ab9f84bc214160efb"
const MELANGE_PRIVATE_KEY = "melange.rsa"
const MELANGE_PUBLIC_KEY = MELANGE_PRIVATE_KEY + ".pub"

const APKO_IMAGE = "cgr.dev/chainguard/apko@sha256:438952dd4da259c6ea728be301e22a85fa21f3bd390a5d66fde19b1ce7c1ba9a"

func GetClient() (string, error) {
	clients := []string{
		"docker",
		"nerdctl",
	}
	for _, client := range clients {
		fmt.Printf("Searching for %s\n", client)
		path, err := exec.LookPath(client)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("Could not find any container clients, please install one of the following %s\n", clients)
}
