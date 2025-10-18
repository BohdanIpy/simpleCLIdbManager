package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(host string, port string, user string, dbname string, password string) (*sql.DB, error) {
	//connString := "host=127.0.0.1 port=5432 user=bohdan sslmode=disable dbname=mydatabase password=12345"
	connString := fmt.Sprintf("host=%s port=%s user=%s sslmode=disable dbname=%s password=%s", host, port, user, dbname, password)
	return sql.Open("postgres", connString)
}
