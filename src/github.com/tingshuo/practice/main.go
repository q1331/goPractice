package main

import (
	"fmt"
	"database/sql"
	"log"
	"github.com/gorilla/mux"
)

func main() {
	a := App{}
	connectionString := fmt.Sprintf("%s:%s@/%s", "root", "sees7&chanting", "practice")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println("Connection string: " + connectionString)
		log.Fatal(err, "Database connection failed.")
	}
	a.Initialize(mux.NewRouter(), db)
	a.initializeRoutes()
	a.Run(":8080")
}

