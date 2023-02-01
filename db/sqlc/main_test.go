package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// setup database
func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load configurations")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to setup database connection")
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
