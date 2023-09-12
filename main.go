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

func writeToFile(content, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func main() {
	home := os.Getenv("HOME")
	sshKey := os.Getenv("PLUGIN_SSH_KEY")
	droneGitSSHURL := os.Getenv("DRONE_GIT_SSH_URL")
	droneCommitBranch := os.Getenv("DRONE_COMMIT_BRANCH")
	droneCommit := os.Getenv("DRONE_COMMIT")

	// Create directories and files
	if err := os.MkdirAll(home+"/.ssh/", 0755); err != nil {
		fmt.Println("Failed to create directory:", err)
		return
	}

	if err := writeToFile(sshKey, home+"/.ssh/id_ed25519"); err != nil {
		fmt.Println("Failed to write SSH key:", err)
		return
	}

	// Set permissions for files
	if err := os.Chmod(home+"/.ssh/id_ed25519", 0600); err != nil {
		fmt.Println("Failed to set permissions for id_ed25519:", err)
		return
	}
	if err := executeCmd("la", home+"/.ssh"); err != nil {
		return
	}
	if err := executeCmd("cat", home+"/.ssh/id_ed25519"); err != nil {
		return
	}
	if _, err := os.Create(home + "/.ssh/known_hosts"); err != nil {
		fmt.Println("Failed to create known_hosts:", err)
		return
	}

	if err := os.Chmod(home+"/.ssh/known_hosts", 0600); err != nil {
		fmt.Println("Failed to set permissions for known_hosts:", err)
		return
	}

	if _, err := os.Create(home + "/.ssh/config"); err != nil {
		fmt.Println("Failed to create config:", err)
		return
	}

	if err := os.Chmod(home+"/.ssh/config", 0600); err != nil {
		fmt.Println("Failed to set permissions for config:", err)
		return
	}

	// Write data to files
	configContent := `Host github.com
    Hostname ssh.github.com
    Port 443
    User git`
	if err := writeToFile(configContent, home+"/.ssh/config"); err != nil {
		fmt.Println("Failed to write to config:", err)
		return
	}

	knownHostsContent := "[ssh.github.com]:443 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl"
	if err := writeToFile(knownHostsContent, home+"/.ssh/known_hosts"); err != nil {
		fmt.Println("Failed to append to known_hosts:", err)
		return
	}

	// Execute git commands
	gitCommands := [][]string{
		{"config", "--global", "init.defaultBranch", "main"},
		{"init"},
		{"config", "advice.detachedHead", "false"},
		{"remote", "add", "origin", droneGitSSHURL},
		{"fetch", "--no-tags", "--prune", "--no-recurse-submodules", "origin", droneCommitBranch},
		{"checkout", droneCommit},
	}

	for _, args := range gitCommands {
		if err := executeCmd("git", args...); err != nil {
			fmt.Printf("Failed to execute git command: git %v\n", args)
			return
		}
	}
}
