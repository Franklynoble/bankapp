package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Franklynoble/bankapp/db/util"
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
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	//var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("can not coonnect to db:", err)

	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
