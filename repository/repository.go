package repository

import (
	"git-996/repository/model"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GitRepository struct {
	Path string
}

func (repo *GitRepository) runCommand(name string, arg ...string) string {
	exeCmd := exec.Command(name, arg...)
	exeCmd.Dir = repo.Path
	output, err := exeCmd.Output()
	if err != nil {
		panic(err.Error())
	} else {
		return string(output)
	}
}

func (repo *GitRepository) AllCommitInfo() []*model.GitCommitInfo {
	var res []*model.GitCommitInfo
	allCommitHash := repo.AllCommitHash()
	hashCh := make(chan string)
	gitCommitInfoCh := make(chan *model.GitCommitInfo)
	var wg sync.WaitGroup
	var done sync.WaitGroup

	go func() {
		for _, commitHash := range allCommitHash {
			hashCh <- commitHash
		}
		close(hashCh)
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for {
				if commitHash, result := <-hashCh; result {
					gitCommitInfoCh <- repo.CommitInfo(commitHash)
				} else {
					wg.Done()
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(gitCommitInfoCh)
	}()

	done.Add(1)
	go func() {
		for {
			if commitInfo, result := <-gitCommitInfoCh; result {
				res = append(res, commitInfo)
			} else {
				done.Done()
				break
			}
		}
	}()

	done.Wait()
	return res
}

func (repo *GitRepository) AllCommitHash() []string {
	result := repo.runCommand("git", "rev-list", "--all")
	result = strings.Trim(result, "\n")
	if result == "" {
		// empty repository
		return make([]string, 0)
	} else {
		return strings.Split(result, "\n")
	}
}

func (repo *GitRepository) CommitInfo(commitHash string) *model.GitCommitInfo {
	result := repo.runCommand("git", "show", "--pretty=format:hash=%H%nname=%an%nemail=%ae%ntime=%at%n==split==", "--shortstat", commitHash)
	parts := strings.SplitN(result, "==split==", 2)
	var res = &model.GitCommitInfo{}
	parseGitCommitMetaInfo(parts[0], res)
	parseGitShortstat(parts[1], res)
	return res
}

func (repo *GitRepository) Summary() map[string]*model.SummaryItem {
	var commitSummary = make(map[string]*model.SummaryItem)
	fileCh := make(chan string)
	blameInfoCh := make(chan map[string]*model.SummaryItem)
	var wg sync.WaitGroup
	var done sync.WaitGroup

	go func() {
		result := repo.runCommand("git", "-c", "core.quotepath=off", "ls-files", "--eol")
		if result == "" {
			return
		}
		result = strings.Trim(result, "\n")
		for _, line := range strings.Split(result, "\n") {
			if !strings.HasPrefix(line, "i/-text") {
				fileCh <- strings.Split(line, "\t")[1]
			}
		}
		close(fileCh)
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for {
				if fileName, result := <-fileCh; result {
					blameInfoCh <- repo.FileBlameInfo(fileName, "HEAD")
				} else {
					wg.Done()
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(blameInfoCh)
	}()

	done.Add(1)
	for {
		if subMap, result := <-blameInfoCh; result {
			for email, summaryItem := range subMap {
				item, err := commitSummary[email]
				if !err {
					item = &model.SummaryItem{Email: email, N: 0}
					commitSummary[email] = item
				}
				item.N += summaryItem.N
			}
		} else {
			done.Done()
			break
		}
	}

	done.Wait()
	return commitSummary
}

func (repo *GitRepository) FileBlameInfo(fileName, hash string) map[string]*model.SummaryItem {
	result := repo.runCommand("git", "blame", "-e", hash, "--", fileName)
	summary := make(map[string]*model.SummaryItem)
	for _, line := range strings.Split(result, "\n") {
		if line != "" {
			email := line[strings.Index(line, "(<")+2:]
			email = email[:strings.Index(email, ">")]
			item, err := summary[email]
			if !err {
				item = &model.SummaryItem{Email: email, N: 0}
				summary[email] = item
			}
			item.N++
		}
	}
	return summary
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
		case "time":
			unixTime, _ := strconv.ParseInt(parts[1], 10, 32)
			target.AuthorTime = time.Unix(unixTime, 0)
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
