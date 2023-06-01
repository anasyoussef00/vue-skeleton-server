package db

import (
	"github.com/jmoiron/sqlx"
	"os"
	"path"
	"strings"
)

var DB *sqlx.DB

func Connect(driver, DSN string) error {
	var err error
	DB, err = sqlx.Connect(driver, DSN)
	return err
}

func Migrate(sqlDirectory string) error {
	entries, err := os.ReadDir(sqlDirectory)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := path.Join(sqlDirectory, entry.Name())
		correctedFullPath := strings.Replace(fullPath, "/", "\\", len(fullPath))

		content, err := os.ReadFile(correctedFullPath)
		if err != nil {
			return err
		}

		DB.MustExec(string(content))
	}
	return nil
}
