package backend_masterclass

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	//We used the lib/pq library to use the correct postgres driver
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

//The convention is to use TestMain as the entry point of all the tests in the application
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(testDB)

	//m.Run() will return exit code that will tell us wether the test passes or fails
	//os.Exit() returns the code to the running program
	os.Exit(m.Run())
}
