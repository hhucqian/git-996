package model

import (
	"encoding/json"
	"fmt"
	"git-996/cmd/arg"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type RepositoryResult_MemberItem struct {
	EMail        string
	Names        map[string]bool
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}

func (item *RepositoryResult_MemberItem) NamesString() string {
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
	Members      map[string]*RepositoryResult_MemberItem
}

type repositoryResultArray struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
	Members      []*RepositoryResult_MemberItem
}

func (repoResult *RepositoryResult) toRepositoryResultArray() repositoryResultArray {
	var rra repositoryResultArray
	rra.CodeDecrease = repoResult.CodeDecrease
	rra.CodeIncrease = repoResult.CodeIncrease
	rra.N = repoResult.N
	for _, v := range repoResult.Members {
		rra.Members = append(rra.Members, v)
	}
	sort.Slice(rra.Members, func(i, j int) bool {
		if arg.RootArg.Revert {
			if strings.HasPrefix(arg.RootArg.Sort, "l") {
				return rra.Members[i].N > rra.Members[j].N
			}
			if strings.HasPrefix(arg.RootArg.Sort, "i") {
				return rra.Members[i].CodeIncrease > rra.Members[j].CodeIncrease
			}
			if strings.HasPrefix(arg.RootArg.Sort, "d") {
				return rra.Members[i].CodeDecrease > rra.Members[j].CodeDecrease
			}
		} else {
			if strings.HasPrefix(arg.RootArg.Sort, "l") {
				return rra.Members[i].N < rra.Members[j].N
			}
			if strings.HasPrefix(arg.RootArg.Sort, "i") {
				return rra.Members[i].CodeIncrease < rra.Members[j].CodeIncrease
			}
			if strings.HasPrefix(arg.RootArg.Sort, "d") {
				return rra.Members[i].CodeDecrease < rra.Members[j].CodeDecrease
			}
		}
		return true
	})
	return rra
}

func (repositoryResult *RepositoryResult) PrintTable() {
	rra := repositoryResult.toRepositoryResultArray()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"e-mail", "name", "+", "-", "产出", "留存率", "贡献率"})
	for _, item := range rra.Members {
		t.AppendRow(table.Row{
			item.EMail,
			item.NamesString(),
			item.CodeIncrease,
			item.CodeDecrease,
			item.N,
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(item.CodeIncrease)*100),
			fmt.Sprintf("%.2f%%", float32(item.N)/float32(repositoryResult.N)*100),
		})
	}
	t.AppendFooter(table.Row{"", "", repositoryResult.CodeIncrease, repositoryResult.CodeDecrease, repositoryResult.N, "-", "-"})
	t.Render()
}

func (repositoryResult *RepositoryResult) PrintJSON() {
	rra := repositoryResult.toRepositoryResultArray()
	res, _ := json.MarshalIndent(rra, "", "  ")
	fmt.Println(string(res))
}
