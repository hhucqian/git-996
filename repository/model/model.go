package model

type GitCommitInfo struct {
	Hash  string
	Name  string
	Email string
	Plus  int32
	Minus int32
}

type GitBlameItem struct {
	Email string
	Hash  string
	N     int32
}
