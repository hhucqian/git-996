package model

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

type GitCommitInfo struct {
	Hash       string
	Name       string
	Email      string
	AuthorTime time.Time
	Plus       int32
	Minus      int32
}

type GitBlameItem struct {
	Email string
	Hash  string
	N     int32
}

type DBCommitSummary struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}

type PrintInfo_UserInfo struct {
	EMail        string
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}
type PrintInfo struct {
	N       int32
	Members map[string]*PrintInfo_UserInfo
}

func (self *PrintInfo) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "+", "-", "产出", "留存率", "贡献率"})
	for email, item := range self.Members {
		t.AppendRow(table.Row{
			email,
			item.CodeIncrease,
			item.CodeIncrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(self.N)*100),
		})
	}
	t.AppendFooter(table.Row{"", "-", "-", self.N, "-", "-"})
	t.Render()
}
