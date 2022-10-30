package get_repository

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
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

type UserCommitAnalyseItem struct {
	Email string
	Name  map[string]bool
	Plus  int
	Minus int
}

func newUserCommitAnalyseItem(email string) *UserCommitAnalyseItem {
	return &UserCommitAnalyseItem{
		Email: email,
		Name:  make(map[string]bool),
		Minus: 0,
		Plus:  0,
	}
}

func (uai *UserCommitAnalyseItem) AddRecord(name string, plus, minus int) {
	uai.Plus += plus
	uai.Minus += minus
	uai.Name[name] = true
}

func New(path string) *GitRepository {
	res := &GitRepository{
		Path:    path,
		Commits: make([]GitCommitInfo, 0, 100),
	}
	return res
}

func (repo *GitRepository) loadCommits() {
	repo.Commits = make([]GitCommitInfo, 0, 100)
	all_commit_hash_list := repo.allCommitHash()
	for _, commit_hash := range all_commit_hash_list {
		repo.Commits = append(repo.Commits, repo.getCommitInfo(commit_hash))
	}
}

func (repo *GitRepository) AnalyseCommit() {
	repo.loadCommits()
	res := make(map[string]*UserCommitAnalyseItem)

	for _, commit_info := range repo.Commits {

		uai, err := res[commit_info.Email]

		if !err {
			uai = newUserCommitAnalyseItem(commit_info.Email)
			res[commit_info.Email] = uai
		}

		for _, record := range commit_info.Diffs {
			uai.AddRecord(commit_info.Name, record.Plus, record.Minus)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "name", "+", "-", "="})

	var total_p, total_m, total_s int

	for _, item := range res {
		names := make([]string, 0)
		for name := range item.Name {
			names = append(names, name)
		}
		total_p += item.Plus
		total_m += item.Minus
		total_s += item.Plus - item.Minus
		t.AppendRow(table.Row{item.Email, strings.Join(names, ","), item.Plus, item.Minus, item.Plus - item.Minus})
		t.AppendSeparator()
	}
	t.AppendFooter(table.Row{"", "", total_p, total_m, total_s})
	t.Render()
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

	if split_line_n == 0 {
		split_line_n = len(lines)
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
