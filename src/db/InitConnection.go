package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Db : database connection object
var db *sql.DB

// InitConnection : initialize db connection
func InitConnection(host string, port string, dbname string, user string, password string) {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	_db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		panic(err)
	} else {
		db = _db
	}
}
