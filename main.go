package main

import (
	"fmt"
	"os"
	"os/exec"
)

func executeCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	commands := []struct {
		cmd  string
		args []string
	}{
		{"mkdir", []string{"-p", "$HOME/.ssh/"}},
		{"echo", []string{"-n", "$SSH_KEY", ">", "$HOME/.ssh/id_ed25519"}},
		{"chmod", []string{"600", "$HOME/.ssh/id_ed25519"}},
		{"touch", []string{"$HOME/.ssh/known_hosts"}},
		{"chmod", []string{"600", "$HOME/.ssh/known_hosts"}},
		{"touch", []string{"$HOME/.ssh/config"}},
		{"chmod", []string{"600", "$HOME/.ssh/config"}},
		{"echo", []string{"Host github.com\n    Hostname ssh.github.com\n    Port 443\n    User git", ">", "$HOME/.ssh/config"}},
		{"echo", []string{"[ssh.github.com]:443 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl", ">>", "$HOME/.ssh/known_hosts"}},
		{"git", []string{"config", "--global", "init.defaultBranch", "main"}},
		{"git", []string{"init"}},
		{"git", []string{"config", "advice.detachedHead", "false"}},
		{"git", []string{"remote", "add", "origin", "$DRONE_GIT_SSH_URL"}},
		{"git", []string{"fetch", "--no-tags", "--prune", "--no-recurse-submodules", "origin", "$DRONE_COMMIT_BRANCH"}},
		{"git", []string{"checkout", "$DRONE_COMMIT"}},
		{"echo", []string{"pwd"}},
		{"ls", []string{"-a"}},
	}

	for _, c := range commands {
		if err := executeCmd(c.cmd, c.args...); err != nil {
			fmt.Printf("Failed to execute command: %s %v\n", c.cmd, c.args)
			return
		}
	}
}
