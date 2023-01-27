package print

import (
	"fmt"
	"os"

	"git-analyse/model"
	"git-analyse/persisit"

	"github.com/jedib0t/go-pretty/table"
)

func Print(dbPath string) {
	var storage *persisit.Storage
	var err error
	storage, err = persisit.NewPersist(dbPath)
	if err != nil {
		panic(err.Error())
	}
	if err := storage.Open(); err != nil {
		panic(err.Error())
	}
	defer storage.Close()

	emails, _ := storage.EmailList()
	fmt.Printf("emails: %v\n", emails)

	lastCommitHash, _ := storage.LatestCommitHash()
	fmt.Printf("lastCommitHash: %v\n", lastCommitHash)

	commitSummary, _ := storage.CommitSummary(lastCommitHash)
	fmt.Printf("commitSummary: %v\n", commitSummary)

	var userCommitSummary = make(map[string]model.DBCommitSummary)
	for _, email := range emails {
		summary, _ := storage.CommitEmailSummary(lastCommitHash, email)
		userCommitSummary[email] = summary
	}
	fmt.Printf("userCommitSummary: %v\n", userCommitSummary)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "+", "-", "产出", "留存率", "贡献率"})
	for email, item := range userCommitSummary {
		t.AppendRow(table.Row{
			email,
			item.CodeIncrease,
			item.CodeIncrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(commitSummary.N)*100),
		})
	}
	t.AppendFooter(table.Row{"", "-", "-", commitSummary.N, "-", "-"})
	t.Render()
}
