package db


import {
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
}

func CreateDatabase() (*sql.DB, error) {
	serverName := 
}