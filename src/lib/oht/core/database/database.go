package database

import (
	"github.com/boltdb/bolt"
	"log"

	"lib/oht/core/common"
)

func InitializeDatabase() {
	db, err := bolt.Open(common.AbsolutePath(common.DefaultDataDir(), "oht.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database initialized")
	defer db.Close()

}
