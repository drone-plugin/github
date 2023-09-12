package main

import (
	"fmt"
	"os"
	"os/exec"
)

type SSHFile struct {
	Path    string
	Content string
	Perm    os.FileMode
}

const configContent = `Host github.com
Hostname ssh.github.com
Port 443
User git`

const knownHostsContent = "[ssh.github.com]:443 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl"

func executeCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeToFile(content, filename string, perm os.FileMode) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Chmod(perm); err != nil {
		return err
	}
	_, err = f.WriteString(content)
	return err
}

func checkError(action string, err error) bool {
	if err != nil {
		fmt.Printf("Failed to %s: %v\n", action, err)
		return true
	}
	return false
}

func main() {
	home := os.Getenv("HOME")
	sshKey := os.Getenv("PLUGIN_SSH_KEY")
	droneGitSSHURL := os.Getenv("DRONE_GIT_SSH_URL")
	droneCommitBranch := os.Getenv("DRONE_COMMIT_BRANCH")
	droneCommit := os.Getenv("DRONE_COMMIT")
	if sshKey == "" {
		fmt.Println("Missing PLUGIN_SSH_KEY environment variable")
		return
	}
	sshDir := home + "/.ssh/"
	if err := os.MkdirAll(sshDir, 0700); err != nil { // 创建.ssh目录，并设置权限为0700
		fmt.Println("Failed to create .ssh directory:", err)
		return
	}
	sshFiles := []SSHFile{
		{
			Path:    home + "/.ssh/id_ed25519",
			Content: sshKey,
			Perm:    0600,
		},
		{
			Path:    home + "/.ssh/config",
			Content: configContent,
			Perm:    0600,
		},
		{
			Path:    home + "/.ssh/known_hosts",
			Content: knownHostsContent,
			Perm:    0600,
		},
	}
	for _, file := range sshFiles {
		if checkError("write to file", writeToFile(file.Content, file.Path, file.Perm)) {
			return
		}
	}

	gitCommands := [][]string{
		{"clone", "--branch", droneCommitBranch, droneGitSSHURL, "."}, // 克隆指定分支到当前目录
		{"checkout", droneCommit},
	}

	for _, args := range gitCommands {
		if checkError(fmt.Sprintf("execute git command: git %v", args), executeCmd("git", args...)) {
			return
		}
	}
}
