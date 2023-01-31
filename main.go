package main

import (
	"git-analyse/model"
	"git-analyse/repository"

	"os"
)

func main() {
	if len(os.Args) == 2 {
		loadAndPrintFromPath(os.Args[1])
	}
}

func loadAndPrintFromPath(repositoryPath string) {
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

	commitSummary := git.CommitSummary(git.CurrentHeadHash())
	for email, gitBlameItem := range commitSummary {
		printInfo.Members[email].N += gitBlameItem.N
		printInfo.N += gitBlameItem.N
	}
	printInfo.Print()
}
