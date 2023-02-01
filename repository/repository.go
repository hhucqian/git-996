package repository

import (
	"fmt"
	"git-996/repository/model"
	"os/exec"
	"strconv"
	"strings"
)

type GitRepository struct {
	Path string
}

func (repo *GitRepository) runCommad(name string, arg ...string) string {
	exe_cmd := exec.Command(name, arg...)
	exe_cmd.Dir = repo.Path
	output, err := exe_cmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(output)
	}
}

func (repo *GitRepository) AllCommitInfo() []*model.GitCommitInfo {
	var res []*model.GitCommitInfo
	allCommitHash := repo.AllCommitHash()
	for _, commitHash := range allCommitHash {
		res = append(res, repo.CommitInfo(commitHash))
	}
	return res
}

func (repo *GitRepository) AllCommitHash() []string {
	result := repo.runCommad("git", "rev-list", "--all")
	result = strings.Trim(result, "\n")
	return strings.Split(result, "\n")
}

func (repo *GitRepository) CommitInfo(commitHash string) *model.GitCommitInfo {
	result := repo.runCommad("git", "show", "--pretty=format:hash=%H%nname=%an%nemail=%ae%n==split==", "--shortstat", commitHash)
	parts := strings.SplitN(result, "==split==", 2)
	var res = &model.GitCommitInfo{}
	parseGitCommitMetaInfo(parts[0], res)
	parseGitShortstat(parts[1], res)
	return res
}

func (repo *GitRepository) Summary() map[string]*model.GitBlameItem {
	var commitSummary = make(map[string]*model.GitBlameItem)
	result := repo.runCommad("git", "-c", "core.quotepath=off", "ls-files", "--eol")
	result = strings.Trim(result, "\n")
	for _, line := range strings.Split(result, "\n") {
		if !strings.HasPrefix(line, "i/-text") {
			repo.FileBlameInfo(strings.Split(line, "\t")[1], "HEAD", commitSummary)
		} else {
			fmt.Printf("line: %v\n", line)
		}
	}
	return commitSummary
}

func (repo *GitRepository) FileBlameInfo(fileName, hash string, summary map[string]*model.GitBlameItem) {
	result := repo.runCommad("git", "blame", "-e", hash, "--", fileName)
	for _, line := range strings.Split(result, "\n") {
		if line != "" {
			email := line[strings.Index(line, "(<")+2:]
			email = email[:strings.Index(email, ">")]
			item, err := summary[email]
			if !err {
				item = &model.GitBlameItem{Email: email, N: 0, Hash: hash}
				summary[email] = item
			}
			item.N++
		}
	}
}

func parseGitCommitMetaInfo(src string, target *model.GitCommitInfo) {
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		switch parts[0] {
		case "hash":
			target.Hash = parts[1]
		case "name":
			target.Name = parts[1]
		case "email":
			target.Email = parts[1]
		}
	}
}

func parseGitShortstat(src string, target *model.GitCommitInfo) {
	src = strings.Trim(src, "\n")
	parts := strings.Split(src, ",")
	for _, part := range parts {
		if strings.HasSuffix(part, "(+)") {
			value, _ := strconv.ParseInt(strings.Split(part, " ")[1], 10, 32)
			target.Plus += int32(value)
		}
		if strings.HasSuffix(part, "(-)") {
			value, _ := strconv.ParseInt(strings.Split(part, " ")[1], 10, 32)
			target.Minus += int32(value)
		}
	}
}
