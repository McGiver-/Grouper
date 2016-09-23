package handlers

import(
  "fmt"
  "github.com/McGiver-/Grouper/src/myDB"
  "encoding/json"
  "crypto/sha512"
  "encoding/base64"
  "net/http"
)

type AddUserResponse struct{
	Success bool    `json:"success"`
	DBError bool    `json:dberror`
	UserId 	int     `json:"userid"`
}

func login(rw http.ResponseWriter, req *http.Request){
  fmt.Println("Login visited")
  db := myDB.Db

  username := req.Url.Query().get("username")
  password := req.Url.Query().get("password")

  hasher := sha512.New()
  password = password +"George is the worst damn coder ever"
  hasher.Write([]byte(password))
  hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

  rows, err := db.Query("SELECT id FROM accounts WHERE username=$1 and password=$2",username,password)

  if err != nil {
		fmt.Println("DATABASE ERROR: Failed username search in adduser")
		rw.WriteHeader(http.StatusConflict)
		if encoded := jsonResponse(&rw,LoginResponse{false,false,-1}); encoded != true{
			return
		}
	}
  var foundId int
  for rows.Next(){
    if err := rows.Scan(&foundId); err != nil {
      fmt.Println("DATABASE ERROR: Failed to scan")
      rw.WriteHeader(http.StatusConflict)
      if encoded := jsonResponse(&rw, LoginResponse{false;false;-1}); encoded != true {
        return
      }
    }
    rw.WriteHeader(http.StatusConflict)
    if encoded := jsonResponse(&rw, AddUserResponse{true,false;foundId}); encoded != true {
      return
    }
  }
}

func jsonResponse(rw *http.ResponseWriter , response AddUserResponse) (bool){

	(*rw).Header().Set("Content-Type","application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil{
		fmt.Println("Error during json encoding in login")
		return false
	}else{
		fmt.Println("Json sent successfully in login")
		return true
	}
}
