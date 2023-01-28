package print

import (
	"git-analyse/persisit"
)

func Print(dbPath string) {
	var storage *persisit.Storage
	var err error
	storage, err = persisit.NewPersist(dbPath)
	if err != nil {
		panic(err.Error())
	}
	if err := storage.Open(); err != nil {
		panic(err.Error())
	}
	defer storage.Close()

	cs := storage.LatestCommitSummary()
	cs.Print()
}
