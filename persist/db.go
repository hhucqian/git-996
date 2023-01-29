package persist

import (
	"database/sql"
	gitModel "git-analyse/git/model"
	persistModel "git-analyse/persist/model"

	"github.com/blockloop/scan"
)

func (storage *Storage) AddCommitItem(src gitModel.GitCommitInfo) (bool, error) {
	if row, err := storage.DB.Query("select * from commit_item where hash = ?", src.Hash); err != nil {
		return false, err
	} else {
		defer row.Close()
		if row.Next() {
			return false, nil
		} else {
			if _, err := storage.DB.Exec("insert into commit_item(hash, email, name, code_increase, code_decrease, time) values(?,?,?,?,?,?)",
				src.Hash, src.Email, src.Name, src.Plus, src.Minus, src.AuthorTime); err != nil {
				return false, err
			} else {
				return true, nil
			}
		}
	}
}

func (storage *Storage) AddCommitSummary(src gitModel.GitBlameItem) (bool, error) {
	if _, err := storage.DB.Exec("insert into commit_summary(hash, email, n) values(?,?,?)", src.Hash, src.Email, src.N); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (storage *Storage) EmailList() ([]string, error) {
	if rows, err := storage.DB.Query("select distinct email from commit_item;"); err != nil {
		return nil, err
	} else {
		var res []string
		if err := scan.Rows(&res, rows); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
}

func (storage *Storage) LatestCommitHash() (string, error) {
	if rows, err := storage.DB.Query("select hash from commit_item order by time desc limit 1"); err != nil {
		return "", err
	} else {
		var res string
		if err := scan.Row(&res, rows); err != nil {
			return "", err
		} else {
			return res, nil
		}
	}
}

func (storage *Storage) CommitSummary(hash string) (persistModel.DBCommitSummary, error) {
	querySql := `
select sum(code_increase) as codeIncrease, sum(code_decrease) as codeDecrease ,
(select sum(n) from commit_summary where hash = ?) as n
from commit_item
where time <= (select time from commit_item where hash = ?)
	`
	var res persistModel.DBCommitSummary
	if rows, err := storage.DB.Query(querySql, hash, hash); err != nil {
		return res, err
	} else {
		if err := scan.Row(&res, rows); err != nil {
			return res, err
		} else {
			return res, nil
		}
	}
}

func (storage *Storage) CommitEmailSummary(hash, email string) (persistModel.DBCommitSummary, error) {
	querySql := `
select sum(code_increase) as codeIncrease, sum(code_decrease) as codeDecrease,
(select sum(n) from commit_summary where hash = ? and email = ?) as n
from commit_item
where time <= (select time from commit_item where hash = ?)
and email = ?
	`
	var res persistModel.DBCommitSummary
	if rows, err := storage.DB.Query(querySql, hash, email, hash, email); err != nil {
		return res, err
	} else {
		if err := scan.Row(&res, rows); err != nil {
			return res, err
		} else {
			return res, nil
		}
	}
}

func (storage *Storage) Close() {
	storage.DB.Close()
	storage.DB = nil
}

func (storage *Storage) Open() error {
	if db, err := sql.Open("sqlite", storage.DbPath+"?_pragma=busy_timeout(1000)"); err != nil {
		return err
	} else {
		storage.DB = db
		return nil
	}
}

func (storage *Storage) initDBTable() error {
	if err := storage.Open(); err != nil {
		return err
	}
	defer storage.Close()
	var init_sql = `
CREATE TABLE commit_item(
	id INTEGER PRIMARY KEY,
	hash TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL,
	name TEXT NOT NULL,
	code_increase INTEGER NOT NULL,
	code_decrease INTEGER NOT NULL,
	time TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_ci_hash ON commit_item (
    hash
);
CREATE TABLE commit_summary(
	id INTEGER PRIMARY KEY,
	hash TEXT NOT NULL,
	email TEXT NOT NULL,
	n INTEGER NOT NULL
);
CREATE INDEX idx_cs_hash ON commit_summary (
    hash
);
`
	if _, err := storage.DB.Exec(init_sql); err != nil {
		return err
	}
	return nil
}
