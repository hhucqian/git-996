package model

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

type PrintInfo_UserInfo struct {
	EMail        string
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}
type PrintInfo struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Members      map[string]*PrintInfo_UserInfo
}

func (printInfo *PrintInfo) Print() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "+", "-", "产出", "留存率", "贡献率"})
	for email, item := range printInfo.Members {
		t.AppendRow(table.Row{
			email,
			item.CodeIncrease,
			item.CodeIncrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(printInfo.N)*100),
		})
	}
	t.AppendFooter(table.Row{"", printInfo.CodeIncrease, printInfo.CodeDecrease, printInfo.N, "-", "-"})
	t.Render()
}
