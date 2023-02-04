package execute

import (
	"git-996/arg"
	"git-996/model"
	"git-996/repository"
	"strings"
)

func LoadAndPrintFromPath(repositoryPath string) {
	var git = repository.GitRepository{Path: repositoryPath}
	var printInfo = model.PrintInfo{
		Members: make(map[string]*model.PrintInfo_MemberItem),
	}

	allCommitInfo := git.AllCommitInfo()
	for _, commitInfo := range allCommitInfo {
		if printInfo.Members[commitInfo.Email] == nil {
			printInfo.Members[commitInfo.Email] = &model.PrintInfo_MemberItem{
				EMail: commitInfo.Email,
				Names: make(map[string]bool),
			}
		}
		printInfo.Members[commitInfo.Email].Names[commitInfo.Name] = true
		printInfo.Members[commitInfo.Email].CodeIncrease += commitInfo.Plus
		printInfo.Members[commitInfo.Email].CodeDecrease += commitInfo.Minus
		printInfo.CodeIncrease += commitInfo.Plus
		printInfo.CodeDecrease += commitInfo.Minus
	}

	commitSummary := git.Summary()
	for email, gitBlameItem := range commitSummary {
		printInfo.Members[email].N += gitBlameItem.N
		printInfo.N += gitBlameItem.N
	}

	for _, mergeInfo := range arg.RootArg.MergeEMail {
		mails := strings.Split(mergeInfo, "=")
		fromMail := mails[1]
		toMail := mails[0]
		if from, ok := printInfo.Members[fromMail]; ok {
			if to, ok := printInfo.Members[toMail]; ok {
				to.CodeDecrease += from.CodeDecrease
				to.CodeIncrease += from.CodeIncrease
				to.N += from.N
				for name := range from.Names {
					to.Names[name] = true
				}
				delete(printInfo.Members, fromMail)
			}
		}
	}
	if arg.RootArg.Format == "table" {
		printInfo.PrintTable()
	}
	if arg.RootArg.Format == "json" {
		printInfo.PrintJSON()
	}
}
