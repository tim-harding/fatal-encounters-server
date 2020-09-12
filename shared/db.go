package shared

import (
	"database/sql"
	"log"

	// Import for postgres driver
	_ "github.com/lib/pq"
)

// Db is the global database connection
var Db *sql.DB

const connectString = `
	host=localhost 
	port=5432 
	user=postgres 
	password=postgres 
	dbname=fatal_encounters 
	sslmode=disable
`

func init() {
	if db, err := sql.Open("postgres", connectString); err == nil {
		Db = db
	} else {
		log.Fatal(err)
	}

	err := Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}
