package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).MarginBottom(1)
	cleanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	dirtyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	bulletStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
)

type DirtyInfo struct {
	Path   string
	Reason string
}

type RepoStatus struct {
	Path   string
	Reason string
	Style  lipgloss.Style
}

func GetVCSInfos() []RepoStatus {
	all_repos := make([]RepoStatus, 0, 4)

	cwd, err := os.Getwd()
	if err != nil {
		return all_repos
	}

	if isGitRepo(cwd) {
		info := isDirty(cwd)
		style := cleanStyle
		if info.Reason != "" {
			style = dirtyStyle
		}
		all_repos = append(all_repos, RepoStatus{Path: info.Path, Reason: info.Reason, Style: style})
	}

	folders, err := os.ReadDir(cwd)
	if err != nil {
		return all_repos
	}

	for _, folder := range folders {
		folderPath := filepath.Join(cwd, folder.Name())

		if isGitRepo(folderPath) {
			info := isDirty(folderPath)
			style := cleanStyle
			if info.Reason != "" {
				style = dirtyStyle
			}
			all_repos = append(all_repos, RepoStatus{Path: info.Path, Reason: info.Reason, Style: style})
		}
	}

	return all_repos
}

func isGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && info.IsDir()
}

func isDirty(repoPath string) DirtyInfo {
	status := exec.Command("git", "status", "--porcelain")
	log := exec.Command("git", "log", "--oneline", "--branches", "--not", "--remotes")

	status.Dir = repoPath
	log.Dir = repoPath

	outStatus, errStatus := status.Output()
	outLog, errLog := log.Output()

	var reasons []string
	if errStatus == nil && len(outStatus) > 0 {
		reasons = append(reasons, "uncommitted changes")
	}
	if errLog == nil && len(outLog) > 0 {
		reasons = append(reasons, "unpushed commits")
	}

	reason := strings.Join(reasons, " and ")

	return DirtyInfo{Path: repoPath, Reason: reason}
}

func Run(cmd *cobra.Command, args []string) error {
	repos := GetVCSInfos()

	home, _ := os.UserHomeDir()

	fmt.Println(headerStyle.Render("Repository Status"))

	for _, info := range repos {
		displayPath := info.Path
		if home != "" && strings.HasPrefix(info.Path, home) {
			displayPath = "~" + strings.TrimPrefix(info.Path, home)
		}

		status := "clean"
		if info.Reason != "" {
			status = info.Reason
		}

		fmt.Println(bulletStyle.Render("•"), info.Style.Render(fmt.Sprintf("%s (%s)", displayPath, status)))
	}

	return nil
}
