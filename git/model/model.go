package model

import "time"

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
