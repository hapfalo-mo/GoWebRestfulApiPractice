package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	db_user := "root"
	db_password := "HJ10xugb123*"
	db_host := "localhost"
	db_name := "golangtestdatabase"

	// Connecting line to mysql database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db_user, db_password, db_host, db_name))
	if err != nil {
		log.Fatal(err)
	}

	// Chec if the connection is working
	err = db.Ping()
	if err != nil {
		log.Fatal("Đã có lỗi xảy ra khi kết nối: %v", err)
	}
	DB = db
}
