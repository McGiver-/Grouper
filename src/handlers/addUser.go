package handlers
import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"database/sql"
	"crypto/sha512"
	"encoding/base64"
	_"github.com/lib/pq"
)



type AddUserResponse struct{
	Success bool `json:"success"`
	Exists bool `json:"exists"`
	DBError bool `json:dberror`
}

func AddUser(rw http.ResponseWriter, req *http.Request){
	fmt.Println("addUser visited")
	db, err := sql.Open("postgres", "postgresql://george@localhost:26257/grouper?sslmode=disable")

	if err != nil {
		fmt.Println("DATABASE ERROR: Failed to connect to datbase in addUser")
		rw.WriteHeader(http.StatusConflict)
		rw.Header().Set("Content-Type","application/json; charset=UTF-8")
		response := AddUserResponse{false,false,true}
		if encoded := jsonResponse(&rw,response); encoded != true{
			return
		}
	}else{
		fmt.Println("Connected to database")
	}

	username := req.URL.Query().Get("username")
	password := req.URL.Query().Get("password")

	rows, err := db.Query("SELECT EXISTS (SELECT 1 FROM accounts WHERE username=$1 LIMIT 1);",username)

	if err != nil {
		fmt.Println("DATABASE ERROR: Failed username search in adduser")
		response := AddUserResponse{false,false,true}
		if encoded := jsonResponse(&rw,response); encoded != true{
			return
		}
	}

	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	defer rows.Close()
	var foundUsername string
	for rows.Next(){
		if err := rows.Scan(&foundUsername); err != nil{
			fmt.Println("DATABASE ERROR: Failed to scan")
			response := AddUserResponse{false,false,true}
			if encoded := jsonResponse(&rw,response); encoded != true{
				return
			}
		}
		fmt.Printf("Found username = %s",foundUsername)
		switch foundUsername {

		case "false":
			if _, err := db.Exec(
				"INSERT INTO accounts (username, password) VALUES ('"+username+"','"+hashedPass+"')"); err != nil {
				log.Fatal(err)
			}else{
				fmt.Printf("Users %s added",username);
				response := AddUserResponse{true,false,false}
				if encoded := jsonResponse(&rw,response); encoded != true{
					return
				}
			}
		case "true":
			response := AddUserResponse{false,true,false}
			fmt.Printf("Username %s already exists",username)
			if encoded := jsonResponse(&rw,response); encoded != true{
				return
			}
		}
	}
}

func jsonResponse(rw *http.ResponseWriter , response AddUserResponse) (bool){
	if err := json.NewEncoder(rw).Encode(response); err != nil{
		(*rw).WriteHeader(http.StatusConflict)
		fmt.Println("Error during json encoding in addUser")
		return false
	}else{
		(*rw).WriteHeader(http.StatusOK)
		(*rw).Header().Set("Content-Type","application/json; charset=UTF-8")
		fmt.Println("Json sent successfully in addUser")
		return true
	}
}