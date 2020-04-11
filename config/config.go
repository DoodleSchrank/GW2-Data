package config

import("database/sql"
       _ "github.com/go-sql-driver/mysql")

func GetMySQLDB() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "web"
	dbPass := "dummy"
	dbName := "data"
	db, err = sql.Open(dbDriver, dbUser + ":" + dbPass + "@/" + dbName)
	return
}
