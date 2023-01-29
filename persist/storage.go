package persist

import (
	"database/sql"
	"errors"
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
