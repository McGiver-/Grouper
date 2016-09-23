package handlers

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"github.com/McGiver-/Grouper/src/myDB"
	_"github.com/lib/pq"
)

type usernamePassArray struct{
	Username string `json:"username"`
	Password string `json:"password"`
}


func ListUsers(rw http.ResponseWriter, req *http.Request){
	fmt.Println("listUser visited")
	db := myDB.Db
	rows, err := db.Query("SELECT username,password FROM accounts")

	if err != nil {
		fmt.Println("Failed the find")
		log.Fatal(err)
	}

	var username, password string

	response := []usernamePassArray{}

	for rows.Next(){
		if err := rows.Scan(&username,&password); err != nil{
			log.Fatal(err)
		}

		fmt.Println("Users found")

		response = append(response,usernamePassArray{username,password})
	}

	if err := json.NewEncoder(rw).Encode(response); err != nil{
		rw.WriteHeader(http.StatusConflict)
		panic(err)
	}else{
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type","application/json; charset=UTF-8")
		fmt.Println("Json success sent")
	}
}
