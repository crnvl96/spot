package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func checkIfRepoHasUncommitedChanges(repoPath string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	return err == nil && len(out) > 0
}

func checkIfRepoHasUnpushedChanges(repoPath string) bool {
	cmd := exec.Command("git", "log", "--oneline", "--branches", "--not", "--remotes")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	return err == nil && len(out) > 0
}

func checkIfDirty(repoPath string) bool {
	return checkIfRepoHasUncommitedChanges(repoPath) || checkIfRepoHasUnpushedChanges(repoPath)
}

func isGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && info.IsDir()
}

func Run(cmd *cobra.Command, args []string) error {
	dirty_repos := make([]string, 0, 4)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if isGitRepo(cwd) && checkIfDirty(cwd) {
		dirty_repos = append(dirty_repos, filepath.Base(cwd))
	}

	folders, err := os.ReadDir(cwd)
	if err != nil {
		return err
	}

	for _, folder := range folders {
		folderPath := filepath.Join(cwd, folder.Name())

		if isGitRepo(folderPath) && checkIfDirty(folderPath) {
			dirty_repos = append(dirty_repos, folder.Name())
		}

	}

	fmt.Println("Dirty repos:", dirty_repos)

	toggle, err := cmd.Flags().GetBool("toggle")
	if err != nil {
		return err
	}

	if toggle {
		fmt.Println("Hey, its Spot with toggle on!")
		return nil
	}

	fmt.Println("Hey, its Spot with toggle off!")
	return nil
}
