package handlers

import(
  "fmt"
  "github.com/McGiver-/Grouper/src/myDB"
  "encoding/json"
  "crypto/sha512"
  "encoding/base64"
  "net/http"
)

type LoginResponse struct{
	Success bool    `json:"success"`
	DBError bool    `json:"dberror"`
	UserId 	int     `json:"userid"`
}

func Login(rw http.ResponseWriter, req *http.Request){
  fmt.Println("Login visited")
  db := myDB.Db

  username := req.URL.Query().Get("username")
  password := req.URL.Query().Get("password")

  hasher := sha512.New()
  password = password +"George is the worst damn coder ever"
  hasher.Write([]byte(password))
  hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
  fmt.Println("Pass Hashed")
  rows, err := db.Query("SELECT id FROM accounts WHERE username=$1 and password=$2 ;",username,hashedPass)

  if err != nil {
		fmt.Println("DATABASE ERROR: Failed username search in adduser")
		rw.WriteHeader(http.StatusConflict)
		if encoded := sendLoginResponse(&rw,LoginResponse{false,false,-1}); encoded != true{
			return
		}
	}
  fmt.Println("Number 1")
  var foundId int
  for rows.Next(){
    if err := rows.Scan(&foundId); err != nil {
      fmt.Println("DATABASE ERROR: Failed to scan")
      rw.WriteHeader(http.StatusConflict)
      if encoded := sendLoginResponse(&rw, LoginResponse{false,false,-1}); encoded != true {
        return
      }
    }
    rw.WriteHeader(http.StatusOK)
    if encoded := sendLoginResponse(&rw, LoginResponse{true,false,foundId}); encoded != true {
      return
    }
  }
}

func sendLoginResponse(rw *http.ResponseWriter , response LoginResponse) (bool){

	(*rw).Header().Set("Content-Type","application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil{
		fmt.Println("Error during json encoding in login")
		return false
	}else{
		fmt.Println("Json sent successfully in login")
		return true
	}
}
