package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

//conn, err := sql.Open(dbDriver, dbSource)

//
//this main file convention is the main entry point for one unit specific package
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("can not coonnect to db:", err)

	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
