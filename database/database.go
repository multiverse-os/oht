package database

import (
	"log"

	"../common"

	"github.com/boltdb/bolt"
)

func InitializeDatabase() {
	db, err := bolt.Open(common.AbsolutePath(common.DefaultDataDir(), "onionwave.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database initialized")
	defer db.Close()

}
