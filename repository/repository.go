package repository

import (
	"git-analyse/repository/model"
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
	res, err := exe_cmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(res[:len(res)-1])
	}
}

func (repo *GitRepository) AllCommitInfo() []*model.GitCommitInfo {
	result := repo.runCommad("git", "log", "--pretty=tformat:==start==%nname=%an%nemail=%ae", "--shortstat")
	lines := strings.Split(result, "\n")
	lineCount := len(lines)
	currentLine := 0
	res := make([]*model.GitCommitInfo, 0, 100)
	var currentCommitInfo *model.GitCommitInfo
	for {
		if lines[currentLine] == "==start==" {
			currentCommitInfo = &model.GitCommitInfo{}
			res = append(res, currentCommitInfo)
			currentLine++
		}

		for {
			if len(lines[currentLine]) == 0 {
				currentLine++
				break
			}

			parts := strings.SplitN(lines[currentLine], "=", 2)
			switch parts[0] {
			case "name":
				currentCommitInfo.Name = parts[1]
			case "email":
				currentCommitInfo.Email = parts[1]
			}

			currentLine++
		}

		parts := strings.Split(lines[currentLine], ",")
		currentLine++
		for _, part := range parts {
			if strings.HasSuffix(part, "(+)") {
				value, _ := strconv.ParseInt(strings.Split(part, " ")[1], 10, 32)
				currentCommitInfo.Plus += int32(value)
			}
			if strings.HasSuffix(part, "(-)") {
				value, _ := strconv.ParseInt(strings.Split(part, " ")[1], 10, 32)
				currentCommitInfo.Minus += int32(value)
			}
		}

		if lineCount == currentLine {
			break
		}
	}
	return res
}

func (repo *GitRepository) CommitSummary(hash string) map[string]*model.GitBlameItem {
	var commitSummary = make(map[string]*model.GitBlameItem)
	files := repo.AllFilesInCommit(hash)
	for _, fileName := range files {
		repo.FileBlameInfo(fileName, hash, commitSummary)
	}
	return commitSummary
}

func (repo *GitRepository) AllFilesInCommit(hash string) []string {
	result := repo.runCommad("git", "-c", "core.quotepath=off", "ls-tree", "--name-only", "-r", hash)
	return strings.Split(result, "\n")
}

func (repo *GitRepository) CurrentHeadHash() string {
	result := repo.runCommad("git", "rev-parse", "HEAD")
	return result
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
