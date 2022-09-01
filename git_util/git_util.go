package git_util

import (
	"os/exec"
	"strconv"
	"strings"
)

type GitRepository struct {
	Path string
}

type GitCommitDiffItem struct {
	Plus  int
	Minus int
	Path  string
}
type GitCommitInfo struct {
	Name  string
	Email string
	Diffs []GitCommitDiffItem
}

func (repo *GitRepository) RunCommad(name string, arg ...string) string {
	exe_cmd := exec.Command(name, arg...)
	exe_cmd.Dir = repo.Path
	res, err := exe_cmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(res[:])
	}
}

func (repo *GitRepository) AllCommit() []string {
	result := repo.RunCommad("git", "log", "--format=%H")
	return strings.Split(strings.TrimSpace(result), "\n")
}

func (repo *GitRepository) CommitInfo(commit string) GitCommitInfo {
	result := repo.RunCommad("git", "show", "--pretty=%an%n%ae", "--numstat", commit)
	result = strings.TrimSpace(result)
	lines := strings.Split(result, "\n")
	res := GitCommitInfo{
		Name:  "",
		Email: "",
		Diffs: make([]GitCommitDiffItem, len(lines)-3),
	}
	res.Name = lines[0]
	res.Email = lines[1]
	for i := 3; i < len(lines); i++ {
		parts := strings.Split(lines[i], "\t")
		res.Diffs[i-3].Plus, _ = strconv.Atoi(parts[0])
		res.Diffs[i-3].Minus, _ = strconv.Atoi(parts[1])
		res.Diffs[i-3].Path = parts[2]
	}
	return res
}
