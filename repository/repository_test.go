package repository

import (
	"git-996/repository/model"
	"os/exec"
	"testing"
)

func TestGitRepository_AllCommitInfo(t *testing.T) {
	type fields struct {
		Path string
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.GitCommitInfo
	}{
		{
			name:   "init",
			fields: fields{Path: "./repo/"},
			want:   []*model.GitCommitInfo{},
		},
		{
			name:   "empty",
			fields: fields{Path: "./repo/"},
			want:   []*model.GitCommitInfo{{Email: "test@test.com", Name: "test"}},
		},
		{
			name:   "test1",
			fields: fields{Path: "./repo/"},
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
			exec.Command("sh", "./testdata/"+tt.name+".sh").Run()

			repo := &GitRepository{
				Path: tt.fields.Path,
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
	exec.Command("sh", "./testdata/clean.sh").Run()
}
