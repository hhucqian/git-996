package persisit

import (
	"database/sql"
	"errors"
	"git-analyse/model"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type Storage struct {
	DbPath string
	DB     *sql.DB
}

func NewPersist(dbPath string) (*Storage, error) {
	var res *Storage = &Storage{DbPath: dbPath}
	if fi, err := os.Stat(dbPath); err != nil {
		createDBFile(res.DbPath)
		res.initDBTable()
	} else {
		if fi.IsDir() {
			return nil, errors.New("数据库路径不可用")
		}
	}
	return res, nil
}

func createDBFile(dbPath string) {
	if strings.Contains(dbPath, "/") {
		split_n := strings.LastIndex(dbPath, "/")
		dbFolder := dbPath[:split_n]
		os.MkdirAll(dbFolder, 0755)
	}
	os.Create(dbPath)
}

func (storage *Storage) LatestCommitSummary() model.PrintInfo {
	var res = model.PrintInfo{Members: make(map[string]*model.PrintInfo_UserInfo)}
	emails, _ := storage.EmailList()

	lastCommitHash, _ := storage.LatestCommitHash()

	commitSummary, _ := storage.CommitSummary(lastCommitHash)
	res.N = commitSummary.N
	res.CodeDecrease = commitSummary.CodeDecrease
	res.CodeIncrease = commitSummary.CodeIncrease

	for _, email := range emails {
		summary, _ := storage.CommitEmailSummary(lastCommitHash, email)
		res.Members[email] = &model.PrintInfo_UserInfo{EMail: email}
		res.Members[email].N = summary.N
		res.Members[email].CodeDecrease = summary.CodeDecrease
		res.Members[email].CodeIncrease = summary.CodeIncrease
	}
	return res
}
