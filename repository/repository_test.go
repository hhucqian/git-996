package repository

import (
	"git-996/repository/model"
	"os/exec"
	"reflect"
	"testing"
)

func TestGitRepository_AllCommitInfo(t *testing.T) {
	tests := []struct {
		name string
		want []*model.GitCommitInfo
	}{
		{
			name: "init",
			want: []*model.GitCommitInfo{},
		},
		{
			name: "empty",
			want: []*model.GitCommitInfo{{Email: "test@test.com", Name: "test"}},
		},
		{
			name: "binary",
			want: []*model.GitCommitInfo{
				{Email: "test@test.com", Name: "test"},
				{Email: "test@test.com", Name: "test"},
				{Email: "test@test.com", Name: "test"},
			},
		},
		{
			name: "test1",
			want: []*model.GitCommitInfo{
				{Email: "test@test.com", Name: "test", Minus: 2},
				{Email: "test@test.com", Name: "test", Plus: 2, Minus: 1},
				{Email: "test@test.com", Name: "test", Minus: 1, Plus: 1},
				{Email: "test@test.com", Name: "test", Plus: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := t.TempDir()
			exec.Command("sh", "./testdata/"+tt.name+".sh", repoPath).Run()
			repo := &GitRepository{
				Path: repoPath,
			}
			got := repo.AllCommitInfo()

			if len(got) != len(tt.want) {
				t.Fatal("len(got) != len(tt.want)")
			}
			for idx, item := range tt.want {
				if !got[idx].EqualForTest(item) {
					t.Fatal("got != want")
				}
			}

		})
	}
}

func TestGitRepository_Summary(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*model.SummaryItem
	}{
		{
			name: "init",
			want: map[string]*model.SummaryItem{},
		},
		{
			name: "empty",
			want: map[string]*model.SummaryItem{},
		},
		{
			name: "binary",
			want: map[string]*model.SummaryItem{},
		},
		{
			name: "summary",
			want: map[string]*model.SummaryItem{"test@test.com": {Email: "test@test.com", N: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoPath := t.TempDir()
			exec.Command("sh", "./testdata/"+tt.name+".sh", repoPath).Run()
			repo := &GitRepository{
				Path: repoPath,
			}

			if got := repo.Summary(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitRepository.Summary() = %v, want %v", got, tt.want)
			}
		})
	}

	exec.Command("sh", "./testdata/clean.sh").Run()
}
