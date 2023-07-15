package testing

import (
	"context"
	"os/exec"
	"strings"
)

// RepoTopLevel retrieves the repository top level from git using 'git rev-parse --show-toplevel'. This function asuming
// that 'git' command is always available.
func RepoTopLevel() (string, error) {
	cmd := exec.CommandContext(context.Background(), "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	repoTopLevel := strings.ReplaceAll(string(out), "\n", "")
	return repoTopLevel, nil
}
