package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type PrintInfo_MemberItem struct {
	EMail        string
	Names        map[string]bool
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}

func (item *PrintInfo_MemberItem) NamesString() string {
	keys := make([]string, 0, len(item.Names))
	for k := range item.Names {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}

type PrintInfo struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Members      map[string]*PrintInfo_MemberItem
}

func (printInfo *PrintInfo) PrintTable() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "name", "+", "-", "产出", "留存率", "贡献率"})
	for _, item := range printInfo.Members {
		t.AppendRow(table.Row{
			item.EMail,
			item.NamesString(),
			item.CodeIncrease,
			item.CodeDecrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(printInfo.N)*100),
		})
	}
	t.AppendFooter(table.Row{"", "", printInfo.CodeIncrease, printInfo.CodeDecrease, printInfo.N, "-", "-"})
	t.Render()
}

func (printInfo *PrintInfo) PrintJSON() {
	res, _ := json.MarshalIndent(printInfo, "", "  ")
	fmt.Println(string(res))
}
