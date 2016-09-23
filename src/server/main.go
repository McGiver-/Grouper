package main

import (
	"encoding/json"
	"fmt"
	"github.com/McGiver-/Grouper/src/handlers"
	"github.com/McGiver-/Grouper/src/myDb"
	"net/http"
	"log"
)
type Category struct {
	Name   string `json:"name"`
	NbUsrs string `json:"nbusrs"`
}

func main() {
	fmt.Println("Application started")
	myDb.Init("postgres", "postgresql://george@localhost:26257/grouper?sslmode=disable")
	http.HandleFunc("/addUser",handlers.AddUser)
	http.HandleFunc("/listUsers",handlers.ListUsers)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func George(rw http.ResponseWriter, req *http.Request) {

	var category string = req.URL.Query().Get("category")

	categories := []Category{
		{category, "44"},
		{"Hockey", "10"},
	}

	if err := json.NewEncoder(rw).Encode(categories); err != nil {
		rw.WriteHeader(http.StatusConflict)
		panic(err)
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Println("George you have accessed json")
	for _, i := range categories {
		fmt.Println("Name: " + i.Name + " Number: " + i.NbUsrs)
	}
}
