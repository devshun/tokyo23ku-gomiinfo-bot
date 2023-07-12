package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
