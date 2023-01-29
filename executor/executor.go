package executor

import (
	"fmt"
	"git-analyse/executor/model"
	"git-analyse/git"
	"git-analyse/persist"
)

func LoadAndPersistFromPath(repositoryPath, dbPath string) error {
	var repo = git.Repository{Path: repositoryPath}
	var storage *persist.Storage
	var err error
	storage, err = persist.NewPersist(dbPath)
	if err != nil {
		return err
	}
	err = storage.Open()
	if err != nil {
		return err
	}
	defer storage.Close()
	commitHashList := repo.AllCommitHash()
	for _, commitHash := range commitHashList {
		commitItem := repo.CommitInfo(commitHash)
		if added, err := storage.AddCommitItem(commitItem); err != nil {
			return err
		} else {
			if added {
				commitSummary := repo.CommitSummary(commitHash)
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

func LoadAndPrintFromPath(repositoryPath string) {
	var repo = git.Repository{Path: repositoryPath}
	var printInfo = model.PrintInfo{
		Members: make(map[string]*model.PrintInfo_UserInfo),
	}
	commitHashList := repo.AllCommitHash()
	if len(commitHashList) == 0 {
		fmt.Println("git repository is empty")
		return
	}
	for _, commitHash := range commitHashList {
		commitInfo := repo.CommitInfo(commitHash)
		if printInfo.Members[commitInfo.Email] == nil {
			printInfo.Members[commitInfo.Email] = &model.PrintInfo_UserInfo{EMail: commitInfo.Email}
		}
		printInfo.Members[commitInfo.Email].CodeIncrease += commitInfo.Plus
		printInfo.Members[commitInfo.Email].CodeDecrease += commitInfo.Minus
		printInfo.CodeIncrease += commitInfo.Plus
		printInfo.CodeDecrease += commitInfo.Minus
	}
	commitSummary := repo.CommitSummary(commitHashList[0])
	for email, gitBlameItem := range commitSummary {
		printInfo.Members[email].N += gitBlameItem.N
		printInfo.N += gitBlameItem.N
	}
	printInfo.Print()
}

func Print(dbPath string) {
	var storage *persist.Storage
	var err error
	storage, err = persist.NewPersist(dbPath)
	if err != nil {
		panic(err.Error())
	}
	if err := storage.Open(); err != nil {
		panic(err.Error())
	}
	defer storage.Close()

	var res = model.PrintInfo{Members: make(map[string]*model.PrintInfo_UserInfo)}
	emails, _ := storage.EmailList()

	lastCommitHash, _ := storage.LatestCommitHash()

	commitSummary, _ := storage.CommitSummary(lastCommitHash)
	res.N = commitSummary.N
	res.CodeDecrease = commitSummary.CodeDecrease
	res.CodeIncrease = commitSummary.CodeIncrease

	for _, email := range emails {
		summary, _ := storage.CommitEmailSummary(lastCommitHash, email)
		res.Members[email] = &model.PrintInfo_UserInfo{EMail: email}
		res.Members[email].N = summary.N
		res.Members[email].CodeDecrease = summary.CodeDecrease
		res.Members[email].CodeIncrease = summary.CodeIncrease
	}
	res.Print()
}
