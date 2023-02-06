package model

import "time"

type GitCommitInfo struct {
	Hash       string
	Name       string
	Email      string
	Plus       int32
	Minus      int32
	AuthorTime time.Time
}

func (gci *GitCommitInfo) EqualForTest(target *GitCommitInfo) bool {
	return gci.Email == target.Email && gci.Name == target.Name && gci.Minus == target.Minus && gci.Plus == target.Plus
}

type SummaryItem struct {
	Email string
	N     int32
}
