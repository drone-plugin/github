package main

import (
	"fmt"
	"log"
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
	log.Println("SSH key path:", sshPath)
	sshKey := os.Getenv("PLUGIN_SSH_KEY")
	sshKeyPath := fmt.Sprintf("%s/id_ed25519", sshPath)
	log.Printf("SSH key path: %s", sshKeyPath)
	if err := os.WriteFile(sshKeyPath, []byte(sshKey), 0600); err != nil {
		fmt.Println("Failed to write SSH key:", err)
		os.Exit(1)
	}
	fi, err := os.Stat(sshKeyPath)
	if err != nil {
		log.Fatalf("Failed to stat SSH key file: %v", err)
	} else {
		log.Printf("SSH key file perm: %o", fi.Mode().Perm())
	}
	log.Println("SSH key path:", sshPath)

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

	if err := os.MkdirAll(sshPath, 0700); err != nil {
		log.Printf("Failed to create .ssh directory: %v", err)
	} else {
		fi, err := os.Stat(sshPath)
		if err != nil {
			log.Printf("Failed to stat .ssh directory: %v", err)
		} else {
			log.Printf(".ssh directory perm: %o", fi.Mode().Perm())
		}
	}

	if err := os.WriteFile(sshPath, []byte(sshKey), 0600); err != nil {
		log.Printf("Failed to write SSH key: %v", err)

	} else {
		fi, err := os.Stat(sshPath)
		if err != nil {
			log.Printf("Failed to stat SSH key file: %v", err)
		} else {
			log.Printf("SSH key file perm: %o", fi.Mode().Perm())
		}
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

	gitURL := os.Getenv("DRONE_GIT_SSH_URL")
	fmt.Println("Git URL:", gitURL)
	if err := runCommand("git", "remote", "add", "origin", gitURL); err != nil {
		fmt.Println("Failed to add git remote:", err)
		os.Exit(1)
	}

	gitBranch := os.Getenv("DRONE_COMMIT_BRANCH")
	if err := runCommand("git", "fetch", "--no-tags", "--prune", "--no-recurse-submodules", "origin", gitBranch); err != nil {
		fmt.Println("Failed to fetch git:", err)
		os.Exit(1)
	}

	gitCommit := os.Getenv("DRONE_COMMIT")
	if err := runCommand("git", "checkout", gitCommit); err != nil {
		fmt.Println("Failed to checkout git:", err)
		os.Exit(1)
	}
}
