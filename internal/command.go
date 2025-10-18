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

type repoStatus struct {
	Path   string
	Reason string
	Style  lipgloss.Style
}

func run(cmd *cobra.Command, args []string) error {
	targets, err := cmd.Flags().GetStringSlice("target")
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		targets = []string{"."}
	}

	allRepos := getVCSInfos(targets)

	var dirtyRepos []repoStatus

	for _, repo := range allRepos {
		if repo.Reason != "" {
			dirtyRepos = append(dirtyRepos, repo)
		}
	}

	if len(dirtyRepos) == 0 {
		fmt.Println(cleanStyle.Render("All repositories are clean!"))
	} else {
		fmt.Println(headerStyle.Render("Dirty Repositories Found"))

		for _, info := range dirtyRepos {
			fmt.Println(bulletStyle.Render("•"), info.Style.Render(fmt.Sprintf("%s (%s)", info.Path, info.Reason)))
		}
	}

	return nil
}

func getVCSInfos(targets []string) []repoStatus {
	all_repos := make([]repoStatus, 0, 4)

	for _, target := range targets {
		walkDir(target, 0, 2, &all_repos)
	}

	return all_repos
}

func walkDir(path string, currentDepth int, maxDepth int, repos *[]repoStatus) {
	if currentDepth > maxDepth {
		return
	}

	if isGitRepo(path) {
		info := checkRepoStatus(path)
		*repos = append(*repos, info)
		return
	}

	if currentDepth < maxDepth {
		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				walkDir(filepath.Join(path, entry.Name()), currentDepth+1, maxDepth, repos)
			}
		}
	}
}

func isGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && info.IsDir()
}

func checkRepoStatus(repoPath string) repoStatus {
	absPath, _ := filepath.Abs(repoPath)

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

	style := cleanStyle

	if reason != "" {
		style = dirtyStyle
	}

	return repoStatus{Path: absPath, Reason: reason, Style: style}
}
