package main

import (
	"fmt"
	"git-analyse/analyse"
	"git-analyse/git_util"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage : ", os.Args[0], "<repository path>")
		return
	}

	repository_path := os.Args[1]
	fmt.Println("repository path:", repository_path)

	repository := git_util.GitRepository{Path: repository_path}

	res := make(map[string]*analyse.UserAnalyseItem)

	commits := repository.AllCommit()

	for _, commit := range commits {
		commit_info := repository.CommitInfo(commit)
		uai, err := res[commit_info.Email]

		if !err {
			uai = analyse.New(commit_info.Email)
			res[commit_info.Email] = uai
		}

		for _, record := range commit_info.Diffs {
			uai.AddRecord(commit_info.Name, record.Plus, record.Minus)
		}
	}

	for _, item := range res {
		names := make([]string, 0)
		for name := range item.Name {
			names = append(names, name)
		}
		fmt.Println("Email:", item.Email, "Name:", strings.Join(names, ","))
		fmt.Println("\t+", item.Plus, "-", item.Minus, "->", item.Plus-item.Minus)
	}

}
