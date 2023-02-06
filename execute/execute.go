package execute

import (
	"git-996/cmd/arg"
	"git-996/model"
	"git-996/repository"
	"strings"
)

func LoadAndPrintFromPath(repositoryPath string) {
	var git = repository.GitRepository{Path: repositoryPath}
	var repositoryResult = model.RepositoryResult{
		Members: make(map[string]*model.RepositoryResult_MemberItem),
		Days:    make(map[string]bool),
	}

	allCommitInfo := git.AllCommitInfo()
	for _, commitInfo := range allCommitInfo {
		if repositoryResult.Members[commitInfo.Email] == nil {
			repositoryResult.Members[commitInfo.Email] = &model.RepositoryResult_MemberItem{
				EMail: commitInfo.Email,
				Names: make(map[string]bool),
				Days:  make(map[string]bool),
			}
		}
		repositoryResult.Members[commitInfo.Email].Names[commitInfo.Name] = true
		repositoryResult.Members[commitInfo.Email].Days[commitInfo.AuthorTime.Format("2006-01-02")] = true
		repositoryResult.Members[commitInfo.Email].CodeIncrease += commitInfo.Plus
		repositoryResult.Members[commitInfo.Email].CodeDecrease += commitInfo.Minus
		repositoryResult.CodeIncrease += commitInfo.Plus
		repositoryResult.CodeDecrease += commitInfo.Minus
		repositoryResult.Days[commitInfo.AuthorTime.Format("2006-01-02")] = true
	}
	if len(allCommitInfo) > 0 {
		repositoryResult.To = allCommitInfo[0].AuthorTime.Format("2006-01-02")
		repositoryResult.From = allCommitInfo[len(allCommitInfo)-1].AuthorTime.Format("2006-01-02")
	}

	commitSummary := git.Summary()
	for email, gitBlameItem := range commitSummary {
		repositoryResult.Members[email].N += gitBlameItem.N
		repositoryResult.N += gitBlameItem.N
	}

	for _, mergeInfo := range arg.RootArg.MergeEMail {
		mails := strings.Split(mergeInfo, "=")
		fromMail := mails[1]
		toMail := mails[0]
		if from, ok := repositoryResult.Members[fromMail]; ok {
			if to, ok := repositoryResult.Members[toMail]; ok {
				to.CodeDecrease += from.CodeDecrease
				to.CodeIncrease += from.CodeIncrease
				to.N += from.N
				for name := range from.Names {
					to.Names[name] = true
				}
				delete(repositoryResult.Members, fromMail)
			}
		}
	}
	if arg.RootArg.Format == "table" {
		repositoryResult.PrintTable()
	}
	if arg.RootArg.Format == "json" {
		repositoryResult.PrintJSON()
	}
}
