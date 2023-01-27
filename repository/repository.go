package repository

import (
	"git-analyse/model"
	"git-analyse/persisit"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	Path string
}

func LoadFromPath(repositoryPath, dbPath string) error {
	var repo = Repository{Path: repositoryPath}
	var storage *persisit.Storage
	var err error
	storage, err = persisit.NewPersist(dbPath)
	if err != nil {
		return err
	}
	err = storage.Open()
	if err != nil {
		return err
	}
	defer storage.Close()
	commitHashList := repo.allCommitHash()
	for _, commitHash := range commitHashList {
		commitItem := repo.commitInfo(commitHash)
		if added, err := storage.AddCommitItem(commitItem); err != nil {
			return err
		} else {
			if added {
				commitSummary := repo.commitSummary(commitHash)
				for _, summaryItem := range commitSummary {
					if _, err := storage.AddCommitSummary(*summaryItem); err != nil {
						return err
					}
				}
			} else {
				break
			}
		}
	}
	return nil
}

func (repo *Repository) runCommad(name string, arg ...string) string {
	exe_cmd := exec.Command(name, arg...)
	exe_cmd.Dir = repo.Path
	res, err := exe_cmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(res)
	}
}

func (repo *Repository) allCommitHash() []string {
	result := repo.runCommad("git", "log", "--format=%H")
	return strings.Split(strings.TrimSpace(result), "\n")
}

func (repo *Repository) commitInfo(hash string) model.GitCommitInfo {
	result := repo.runCommad("git", "show", "--pretty=%an%n%ae%n%at", "--numstat", hash)
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

	res := model.GitCommitInfo{Hash: hash}
	res.Name = lines[0]
	res.Email = lines[1]
	unix_time, _ := strconv.ParseInt(lines[2], 10, 64)
	res.AuthorTime = time.Unix(unix_time, 0)
	for i := split_line_n; i < len(lines); i++ {
		parts := strings.Split(lines[i], "\t")
		pValue, _ := strconv.Atoi(parts[0])
		mValue, _ := strconv.Atoi(parts[1])
		res.Plus += pValue
		res.Minus += mValue
	}
	return res
}

func (repo *Repository) commitSummary(hash string) map[string]*model.GitBlameItem {
	var commitSummary = make(map[string]*model.GitBlameItem)
	files := repo.allFilesInCommit(hash)
	for _, fileName := range files {
		repo.fileBlameInfo(fileName, hash, commitSummary)
	}
	return commitSummary
}

func (repo *Repository) allFilesInCommit(hash string) []string {
	result := repo.runCommad("git", "-c", "core.quotepath=off", "ls-tree", "--name-only", "-r", hash)
	return strings.Split(strings.TrimSpace(result), "\n")
}

func (repo *Repository) fileBlameInfo(fileName, hash string, summary map[string]*model.GitBlameItem) {
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
