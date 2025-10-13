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

func GetVCSInfos(targets []string, depth int) []RepoStatus {
	all_repos := make([]RepoStatus, 0, 4)

	for _, target := range targets {
		if strings.HasSuffix(target, "/**") {
			base := strings.TrimSuffix(target, "/**")
			collectReposRecursive(base, depth, &all_repos)
		} else {
			if isGitRepo(target) {
				info := isDirty(target)
				style := cleanStyle
				if info.Reason != "" {
					style = dirtyStyle
				}
				all_repos = append(all_repos, RepoStatus{Path: info.Path, Reason: info.Reason, Style: style})
			}
		}
	}

	return all_repos
}

func collectReposRecursive(base string, maxDepth int, repos *[]RepoStatus) {
	walkDir(base, 0, maxDepth, repos)
}

func walkDir(path string, currentDepth int, maxDepth int, repos *[]RepoStatus) {
	if currentDepth > maxDepth {
		return
	}

	if currentDepth > 0 && isGitRepo(path) {
		info := isDirty(path)
		style := cleanStyle
		if info.Reason != "" {
			style = dirtyStyle
		}
		*repos = append(*repos, RepoStatus{Path: info.Path, Reason: info.Reason, Style: style})
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
	depth, err := cmd.Flags().GetInt("depth")
	if err != nil {
		return err
	}

	targets, err := cmd.Flags().GetStringSlice("target")
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		return fmt.Errorf("no targets specified")
	}

	home, _ := os.UserHomeDir()
	for i, target := range targets {
		if strings.HasPrefix(target, "~") {
			targets[i] = strings.Replace(target, "~", home, 1)
		}
	}

	repos := GetVCSInfos(targets, depth)

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
