package repository

import (
	"fmt"
	"git-analyse/model"
	"git-analyse/persisit"
)

func LoadAndPersistFromPath(repositoryPath, dbPath string) error {
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

func LoadAndPrintFromPath(repositoryPath string) {
	var repo = Repository{Path: repositoryPath}
	var printInfo = model.PrintInfo{
		Members: make(map[string]*model.PrintInfo_UserInfo),
	}
	commitHashList := repo.allCommitHash()
	if len(commitHashList) == 0 {
		fmt.Println("git repository is empty")
		return
	}
	for _, commitHash := range commitHashList {
		commitInfo := repo.commitInfo(commitHash)
		if printInfo.Members[commitInfo.Email] == nil {
			printInfo.Members[commitInfo.Email] = &model.PrintInfo_UserInfo{EMail: commitInfo.Email}
		}
		printInfo.Members[commitInfo.Email].CodeIncrease += commitInfo.Plus
		printInfo.Members[commitInfo.Email].CodeDecrease += commitInfo.Minus
		printInfo.CodeIncrease += commitInfo.Plus
		printInfo.CodeDecrease += commitInfo.Minus
	}
	commitSummary := repo.commitSummary(commitHashList[0])
	for email, gitBlameItem := range commitSummary {
		printInfo.Members[email].N += gitBlameItem.N
		printInfo.N += gitBlameItem.N
	}
	printInfo.Print()
}
