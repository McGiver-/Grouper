package handlers

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/McGiver-/Grouper/src/myDB"
	"net/http"
)

type LoginResponse struct {
	Success bool `json:"success"`
	DBError bool `json:"dberror"`
	UserId  int  `json:"userid"`
}

func Login(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Login visited")
	db := myDB.Db

	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	hasher := sha512.New()
	password = password + "George is the worst damn coder ever"
	hasher.Write([]byte(password))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	rows, err := db.Query("SELECT id FROM accounts WHERE username=$1 and password=$2 ;", username, hashedPass)

	if err != nil {
		fmt.Println("DATABASE ERROR: Failed username search in login")
		rw.WriteHeader(http.StatusConflict)
		sendLoginResponse(&rw, LoginResponse{false, false, -1})
	}

	var foundId int

	if rows.Next() {
		rows.Scan(&foundId)
		rw.WriteHeader(http.StatusOK)

		encoded := sendLoginResponse(&rw, LoginResponse{true, false, foundId})
		if encoded != true {
			return
		}
	} else {
		fmt.Println("Incorrect Username and Password")
		rw.WriteHeader(http.StatusConflict)
		sendLoginResponse(&rw, LoginResponse{false, false, -1})
	}
}

func sendLoginResponse(rw *http.ResponseWriter, response LoginResponse) bool {

	(*rw).Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil {
		fmt.Println("Error during json encoding in login")
		return false
	} else {
		fmt.Println("Json sent successfully in login")
		return true
	}
}
