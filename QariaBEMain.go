package main

import (
	"log"
	"flag"
	"net/http"

	"github.com/enricod/qaria-be/db"
	"github.com/enricod/qaria-be/router"
)

func main() {
	port := flag.String("port", "8080", "porta servizio")
	flag.Parse()

	db.DbInit()
	router := router.NewRouter()

	log.Printf("Avvio server su porta %v\n", *port)
	log.Fatal(http.ListenAndServe(":" + *port, router))
}