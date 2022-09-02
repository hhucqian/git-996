package git_util

import (
	"fmt"
	"git-analyse/analyse"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type GitRepository struct {
	Path    string
	Commits []GitCommitInfo
}

type GitCommitDiffItem struct {
	Plus  int
	Minus int
	Path  string
}
type GitCommitInfo struct {
	Name       string
	Email      string
	AuthorTime time.Time
	Diffs      []GitCommitDiffItem
}

func New(path string) *GitRepository {
	res := &GitRepository{
		Path:    path,
		Commits: make([]GitCommitInfo, 0, 100),
	}
	return res
}

func (repo *GitRepository) Load() {
	repo.Commits = make([]GitCommitInfo, 0, 100)
	all_commit_hash_list := repo.allCommitHash()
	for _, commit_hash := range all_commit_hash_list {
		repo.Commits = append(repo.Commits, repo.getCommitInfo(commit_hash))
	}
}

func (repo *GitRepository) Analyse() {
	res := make(map[string]*analyse.UserAnalyseItem)

	for _, commit_info := range repo.Commits {

		uai, err := res[commit_info.Email]

		if !err {
			uai = analyse.New(commit_info.Email)
			res[commit_info.Email] = uai
		}

		for _, record := range commit_info.Diffs {
			uai.AddRecord(commit_info.Name, record.Plus, record.Minus)
		}
	}

	for _, item := range res {
		names := make([]string, 0)
		for name := range item.Name {
			names = append(names, name)
		}
		fmt.Println("Email:", item.Email, "Name:", strings.Join(names, ","))
		fmt.Println("\t+", item.Plus, "-", item.Minus, "->", item.Plus-item.Minus)
	}
}

func (repo *GitRepository) allCommitHash() []string {
	result := repo.runCommad("git", "log", "--format=%H")
	return strings.Split(strings.TrimSpace(result), "\n")
}

func (repo *GitRepository) runCommad(name string, arg ...string) string {
	exe_cmd := exec.Command(name, arg...)
	exe_cmd.Dir = repo.Path
	res, err := exe_cmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(res[:])
	}
}
func (repo *GitRepository) getCommitInfo(commit string) GitCommitInfo {
	result := repo.runCommad("git", "show", "--pretty=%an%n%ae%n%at", "--numstat", commit)
	result = strings.TrimSpace(result)
	lines := strings.Split(result, "\n")
	split_line_n := 0
	for i, v := range lines {
		if len(v) == 0 {
			split_line_n = i + 1
			break
		}
	}
	res := GitCommitInfo{
		Name:  "",
		Email: "",
		Diffs: make([]GitCommitDiffItem, len(lines)-split_line_n),
	}
	res.Name = lines[0]
	res.Email = lines[1]
	unix_time, _ := strconv.ParseInt(lines[2], 10, 64)
	res.AuthorTime = time.Unix(unix_time, 0)
	for i := split_line_n; i < len(lines); i++ {
		parts := strings.Split(lines[i], "\t")
		res.Diffs[i-split_line_n].Plus, _ = strconv.Atoi(parts[0])
		res.Diffs[i-split_line_n].Minus, _ = strconv.Atoi(parts[1])
		res.Diffs[i-split_line_n].Path = parts[2]
	}
	return res
}
