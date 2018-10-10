// Package homedir detects the user's home directory without the use of cgo, for use in cross-compilation environments.

// +build darwin

package homedir

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func dir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != "" {
		return home, nil
	}

	stdout := new(bytes.Buffer)

	cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
	cmd.Stdout = stdout
	if err := cmd.Run(); err == nil {
		result := string(bytes.TrimSpace(stdout.Bytes()))
		if result != "" {
			return result, nil
		}
	}

	// If all else fails, try the shell
	stdout.Reset()

	cmd = exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := string(bytes.TrimSpace(stdout.Bytes()))
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}
