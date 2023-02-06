package model

import (
	"encoding/json"
	"fmt"
	"git-996/cmd/arg"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type RepositoryresultMemberitem struct {
	EMail        string
	Names        map[string]bool
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Days         map[string]bool
}

func (item *RepositoryresultMemberitem) NamesString() string {
	keys := make([]string, 0, len(item.Names))
	for k := range item.Names {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}

type RepositoryResult struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Members      map[string]*RepositoryresultMemberitem
	Days         map[string]bool
	From         string
	To           string
}

type repositoryResultArray struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Members      []*RepositoryresultMemberitem
}

func (repositoryResult *RepositoryResult) toRepositoryResultArray() repositoryResultArray {
	var rra repositoryResultArray
	rra.CodeDecrease = repositoryResult.CodeDecrease
	rra.CodeIncrease = repositoryResult.CodeIncrease
	rra.N = repositoryResult.N
	for _, v := range repositoryResult.Members {
		rra.Members = append(rra.Members, v)
	}
	return rra
}

func (repositoryResult *RepositoryResult) PrintTable() {
	rra := repositoryResult.toRepositoryResultArray()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "name", "days", "+", "-", "产出", "留存率", "贡献率"})
	for _, item := range rra.Members {
		t.AppendRow(table.Row{
			item.EMail,
			item.NamesString(),
			len(item.Days),
			item.CodeIncrease,
			item.CodeDecrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(repositoryResult.N)*100),
		})
	}
	t.AppendFooter(table.Row{
		"From：" + repositoryResult.From,
		"To：" + repositoryResult.To,
		len(repositoryResult.Days),
		repositoryResult.CodeIncrease,
		repositoryResult.CodeDecrease,
		repositoryResult.N,
		fmt.Sprintf("%.2f%%", float32(repositoryResult.N)/float32(repositoryResult.CodeIncrease)*100),
		""})
	sortMode := table.AscNumeric
	if arg.RootArg.Revert {
		sortMode = table.DscNumeric
	}
	if strings.HasPrefix(arg.RootArg.Sort, "l") {
		t.SortBy([]table.SortBy{{Name: "产出", Mode: sortMode}})
	} else if strings.HasPrefix(arg.RootArg.Sort, "i") {
		t.SortBy([]table.SortBy{{Name: "+", Mode: sortMode}})
	} else if strings.HasPrefix(arg.RootArg.Sort, "d") {
		t.SortBy([]table.SortBy{{Name: "-", Mode: sortMode}})
	}
	t.SetStyle(table.StyleColoredDark)
	t.Render()
}

func (repositoryResult *RepositoryResult) PrintJSON() {
	rra := repositoryResult.toRepositoryResultArray()
	res, _ := json.MarshalIndent(rra, "", "  ")
	fmt.Println(string(res))
}
