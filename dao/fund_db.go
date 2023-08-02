package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "root:123456@tcp(localhost:3306)/fund?loc=Local")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

func executeWithTransactional(f func(*sql.Tx) error) {

	conn, err := Db.Begin()
	if err != nil {
		fmt.Println("begin failed :", err)
		return
	}

	err = f(conn)
	if err != nil {
		conn.Rollback()
		return
	}

	conn.Commit()
}
