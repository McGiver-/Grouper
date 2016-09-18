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
	Success bool    `json:"success"`
	Exists  bool    `json:"exists"`
	DBError bool    `json:dberror`
	UserId      int     `json:"userid"`
}

func AddUser(rw http.ResponseWriter, req *http.Request){
	fmt.Println("addUser visited")
	db, err := sql.Open("postgres", "postgresql://george@localhost:26257/grouper?sslmode=disable")

	if err != nil || db.Ping() != nil{
		fmt.Println("DATABASE ERROR: Failed to connect to datbase in addUser")
		rw.WriteHeader(http.StatusConflict)
		rw.Header().Set("Content-Type","application/json; charset=UTF-8")
		if encoded := jsonResponse(&rw,AddUserResponse{false,false,true,-1}); encoded != true{
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
		rw.WriteHeader(http.StatusConflict)
		if encoded := jsonResponse(&rw,AddUserResponse{false,false,true,-1}); encoded != true{
			return
		}
	}

	hasher := sha512.New()
	password = password + "George is the worst damn coder ever"
	hasher.Write([]byte(password))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	defer rows.Close()
	var foundUsername string
	for rows.Next(){
		if err := rows.Scan(&foundUsername); err != nil{
			fmt.Println("DATABASE ERROR: Failed to scan")
			rw.WriteHeader(http.StatusConflict)
			if encoded := jsonResponse(&rw,AddUserResponse{false,false,true,-1}); encoded != true{
				return
			}
		}
		fmt.Printf("Found username = %s \n",foundUsername)
		switch foundUsername {

		case "false":
			if _, err := db.Exec(
				"INSERT INTO accounts (username, password) VALUES ('"+username+"','"+hashedPass+"')"); err != nil {
				log.Fatal(err)
			}else{
				rows,err = db.Exec("SELECT id FROM accounts WHERE username=$1 , password=$2 ;",username, hashedPass)
				if err != nil{
					fmt.Println("DATABASE ERROR: Failed to get id after insert in adduser")
					rw.WriteHeader(http.StatusConflict)
					if encoded := jsonResponse(&rw,AddUserResponse{false,false,true,-1}); encoded != true{
						return
					}
				}else{
					var foundId int
					fmt.Printf("Users %s added \n",username);
					for rows.Next() {
						if err := rows.Scan(&foundId); err != nil {
							fmt.Println("DATABASE ERROR: Failed to scan")
							rw.WriteHeader(http.StatusOK)
							if encoded := jsonResponse(&rw, AddUserResponse{false, false, true, foundId}); encoded != true {
								return
							}
						}
					}
				}
			}
		case "true":
			fmt.Printf("Username %s already exists \n",username)
			rw.WriteHeader(http.StatusConflict)
			if encoded := jsonResponse(&rw,AddUserResponse{false,true,false,-1}); encoded != true{
				return
			}
		}
	}
}

func jsonResponse(rw *http.ResponseWriter , response AddUserResponse) (bool){

	(*rw).Header().Set("Content-Type","application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil{
		fmt.Println("Error during json encoding in addUser")
		return false
	}else{
		fmt.Println("Json sent successfully in addUser")
		return true
	}
}