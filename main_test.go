package main

import (
	"database/sql"
	"log"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

func Test_fetchData(t *testing.T) {
	var wg *sync.WaitGroup
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password= 12345 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Conn Open → ", err)
	}
	defer db.Close()
	t.Run(tt.name, func(t *testing.T) {
		fetchData("USD", wg, db)
	})

}
