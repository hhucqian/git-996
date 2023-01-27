package model

import "time"

type GitCommitInfo struct {
	Hash       string
	Name       string
	Email      string
	AuthorTime time.Time
	Plus       int
	Minus      int
}

type GitBlameItem struct {
	Email string
	Hash  string
	N     int
}

type DBCommitSummary struct {
	CodeIncrease int32
	CodeDecrease int32
	N            int32
}
