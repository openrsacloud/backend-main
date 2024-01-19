package database

import (
	"os"

	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

func Connect() {
	db, err := surrealdb.New("ws://truenas.local:30888/rpc")
	if err != nil {
		panic(err)
	}
	if _, err := db.Signin(map[string]interface{}{
		"ns":   "OpenRSACloud",
		"db":   "main",
		"user": os.Getenv("SurrealDatabaseUser"),
		"pass": os.Getenv("SurrealDatabasePass"),
	}); err != nil {
		panic(err)
	}
	// if _, err := db.Use("OpenRSACloud", "main"); err != nil {
	// 	panic(err)
	// }
	DB = db
}

func Disconnect() {
	DB.Close()
}
