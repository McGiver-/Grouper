package main

import (
	"encoding/json"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"net/http"
	"log"
)

type Category struct {
	Name   string `json:"name"`
	NbUsrs string `json:"nbusrs"`
}

type AddUserResponse struct{
	Success bool `json:"success"`
}

func main() {

//	http.HandleFunc("/categories", George)
//	http.HandleFunc("/friends", FriendHandler)
	http.HandleFunc("/addUser",addUser)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

/* --------------------------------------- Handlers -------------------------------------------------- */
func addUser(rw http.ResponseWriter, req *http.Request){
	db, err := sql.Open("postgres", "postgresql://george@localhost:26257/accounts?sslmode=disable")
	if err != nil {
		log.Fatalf("error connection to the database: %s", err)
	}

	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	rows, err := db.Query("SELECT EXISTS (SELECT 1 FROM accounts WHERE username=? LIMIT 1);",username)

	if err != nil {
		fmt.Println("Failed the find")
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next(){
		if err := rows.Scan(&username); err != nil{
			log.Fatal(err)
		}
		fmt.Println("username: "+ username)
	}

	// Insert two rows into the "accounts" table.
	if _, err := db.Exec(
		"INSERT INTO accounts (username, password) VALUES ("+username+","+password+")"); err != nil {
		log.Fatal(err)
	}else{
		fmt.Printf("Users %s added",username);
		response := AddUserResponse{true}
		if err := json.NewEncoder(rw).Encode(response); err != nil{
			rw.WriteHeader(http.StatusConflict)
			panic(err)
		}else{
			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type","application/json; charset=UTF-8")
			fmt.Println("Json success sent")
		}
	}
}


func FriendHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Got friends")
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