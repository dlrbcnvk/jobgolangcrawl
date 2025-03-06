package database

import (
	"database/sql"
	"jobgolangcrawl/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Initialize(config *config.Config) *sql.DB {
	// 데이터베이스 연결 설정
	dsn := config.DB.Url
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 데이터베이스 연결 확인
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
