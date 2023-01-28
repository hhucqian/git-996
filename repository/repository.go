package repository

import "git-analyse/persisit"

func LoadAndPersistFromPath(repositoryPath, dbPath string) error {
	var repo = Repository{Path: repositoryPath}
	var storage *persisit.Storage
	var err error
	storage, err = persisit.NewPersist(dbPath)
	if err != nil {
		return err
	}
	err = storage.Open()
	if err != nil {
		return err
	}
	defer storage.Close()
	commitHashList := repo.allCommitHash()
	for _, commitHash := range commitHashList {
		commitItem := repo.commitInfo(commitHash)
		if added, err := storage.AddCommitItem(commitItem); err != nil {
			return err
		} else {
			if added {
				commitSummary := repo.commitSummary(commitHash)
				for _, summaryItem := range commitSummary {
					if _, err := storage.AddCommitSummary(*summaryItem); err != nil {
						return err
					}
				}
			} else {
				break
			}
		}
	}
	return nil
}
