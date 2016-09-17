package handlers
import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"database/sql"
	"crypto/sha512"
	"encoding/base64"
)



type AddUserResponse struct{
	Success bool `json:"success"`
}

func AddUser(rw http.ResponseWriter, req *http.Request){

	db, err := sql.Open("postgres", "postgresql://george@localhost:26257/grouper?sslmode=disable")
	if err != nil {
		log.Fatalf("error connection to the database: %s", err)
	}

	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	rows, err := db.Query("SELECT EXISTS (SELECT 1 FROM accounts WHERE username=$1 LIMIT 1);",username)

	if err != nil {
		fmt.Println("Failed the find")
		log.Fatal(err)
	}

	hasher := sha512.New()
	hasher.Write(password)
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	defer rows.Close()

	for rows.Next(){
		if err := rows.Scan(&username); err != nil{
			log.Fatal(err)
		}
		if username != false{
			if _, err := db.Exec(
				"INSERT INTO accounts (username, password) VALUES ('"+username+"','"+hashedPass+")"); err != nil {
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
		}else{
			rw.WriteHeader(http.StatusConflict)
			rw.Header().Set("Content-Type","application/json; charset=UTF-8")
			fmt.Println("User already exists")
		}
	}
}

