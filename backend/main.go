package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// connect with the postgres server
// connect to the cybercare extension server(whitelist/blacklist)
// this is for Mac users only
func main() {
	connectionString := "host=database-1.cl0i0y628wpg.us-east-1.rds.amazonaws.com port=5432 user=postgres password=jjw99025 dbname=cybercare sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}
}
