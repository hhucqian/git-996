package main

import (
	"fmt"
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
		Members: make(map[string]*model.PrintInfo_UserInfo),
	}
	commitHashList := git.AllCommitHash()
	if len(commitHashList) == 0 {
		fmt.Println("git repository is empty")
		return
	}
	for _, commitHash := range commitHashList {
		commitInfo := git.CommitInfo(commitHash)
		if printInfo.Members[commitInfo.Email] == nil {
			printInfo.Members[commitInfo.Email] = &model.PrintInfo_UserInfo{EMail: commitInfo.Email}
		}
		printInfo.Members[commitInfo.Email].CodeIncrease += commitInfo.Plus
		printInfo.Members[commitInfo.Email].CodeDecrease += commitInfo.Minus
		printInfo.CodeIncrease += commitInfo.Plus
		printInfo.CodeDecrease += commitInfo.Minus
	}
	commitSummary := git.CommitSummary(commitHashList[0])
	for email, gitBlameItem := range commitSummary {
		printInfo.Members[email].N += gitBlameItem.N
		printInfo.N += gitBlameItem.N
	}
	printInfo.Print()
}
