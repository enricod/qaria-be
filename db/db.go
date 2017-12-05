package db

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func DbInit() {
	var err error
	Db, err = sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/qaria")

	Db.SetMaxIdleConns(2)
	Db.SetMaxOpenConns(20)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("creazione connessione al DB %v", Db)

	//defer Db.Close()
}
