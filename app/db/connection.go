package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (*sql.DB, error) {
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/%s",
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
