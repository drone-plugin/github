package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	// Setup SSH key
	sshPath := fmt.Sprintf("%s/.ssh", os.Getenv("HOME"))
	if err := os.MkdirAll(sshPath, 0700); err != nil {
		fmt.Println("Failed to create .ssh directory:", err)
		os.Exit(1)
	}

	sshKey := os.Getenv("PLUGIN_SSH_KEY")
	if err := os.WriteFile(fmt.Sprintf("%s/id_ed25519", sshPath), []byte(sshKey), 0600); err != nil {
		fmt.Println("Failed to write SSH key:", err)
		os.Exit(1)
	}

	// ...其他SSH设置...

	configContent := `Host github.com
    Hostname ssh.github.com
    Port 443
    User git`

	if err := os.WriteFile(fmt.Sprintf("%s/config", sshPath), []byte(configContent), 0600); err != nil {
		fmt.Println("Failed to write SSH config:", err)
		os.Exit(1)
	}

	knownHostsContent := "[ssh.github.com]:443 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl"
	if err := os.WriteFile(fmt.Sprintf("%s/known_hosts", sshPath), []byte(knownHostsContent), 0600); err != nil {
		fmt.Println("Failed to write known hosts:", err)
		os.Exit(1)
	}

	// Git Operations
	if err := runCommand("git", "config", "--global", "init.defaultBranch", "main"); err != nil {
		fmt.Println("Failed to set git config:", err)
		os.Exit(1)
	}

	if err := runCommand("git", "init"); err != nil {
		fmt.Println("Failed to init git:", err)
		os.Exit(1)
	}

	if err := runCommand("git", "config", "advice.detachedHead", "false"); err != nil {
		fmt.Println("Failed to set git advice:", err)
		os.Exit(1)
	}

	gitURL := os.Getenv("PLUGIN_DRONE_GIT_SSH_URL")
	if err := runCommand("git", "remote", "add", "origin", gitURL); err != nil {
		fmt.Println("Failed to add git remote:", err)
		os.Exit(1)
	}

	gitBranch := os.Getenv("PLUGIN_DRONE_COMMIT_BRANCH")
	if err := runCommand("git", "fetch", "--no-tags", "--prune", "--no-recurse-submodules", "origin", gitBranch); err != nil {
		fmt.Println("Failed to fetch git:", err)
		os.Exit(1)
	}

	gitCommit := os.Getenv("PLUGIN_DRONE_COMMIT")
	if err := runCommand("git", "checkout", gitCommit); err != nil {
		fmt.Println("Failed to checkout git:", err)
		os.Exit(1)
	}
}
