package backend_masterclass

import (
	"backend_masterclass/util"
	"database/sql"
	"log"
	"os"
	"testing"

	//We used the lib/pq library to use the correct postgres driver
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

//The convention is to use TestMain as the entry point of all the tests in the application
func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config values from file: ", err.Error())
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(testDB)

	//m.Run() will return exit code that will tell us wether the test passes or fails
	//os.Exit() returns the code to the running program
	os.Exit(m.Run())
}
