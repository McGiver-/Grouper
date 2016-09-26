package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/McGiver-/Grouper/src/myDB"
	"net/http"
	"strings"
)

type CreateCategoryResponse struct {
	Success bool `json:"success"`
	Exists  bool `json:"exists"`
}

func CreateCategory(rw http.ResponseWriter, req *http.Request) {
	db := myDB.Db
	category := req.URL.Query().Get("category")
	category = strings.ToLower(category)

	rows, err := db.Query("SELECT name FROM categories WHERE name=$1 ;", category)

	if rows.Next() {
		fmt.Printf("Table %s Exists \n", category)
		rw.WriteHeader(http.StatusConflict)
		sendCreateCategoryResponse(&rw, CreateCategoryResponse{false, true})
		return
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " + category + " (userId BIGSERIAL);")

	if err != nil {
		fmt.Printf("Could not create table %s \n", category)
		rw.WriteHeader(http.StatusConflict)
		sendCreateCategoryResponse(&rw, CreateCategoryResponse{false, false})
		return
	} else {
		_, err := db.Exec("INSERT INTO categories (name) VALUES ('" + category + "') ;")
		if err != nil {
			fmt.Printf("Could not add the category %s to categories table \n", category)
			rw.WriteHeader(http.StatusConflict)
			sendCreateCategoryResponse(&rw, CreateCategoryResponse{false, false})
			return
		} else {
			rw.WriteHeader(http.StatusOK)
			sendCreateCategoryResponse(&rw, CreateCategoryResponse{true, false})
		}
	}
}

func sendCreateCategoryResponse(rw *http.ResponseWriter, response CreateCategoryResponse) bool {

	(*rw).Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(*rw).Encode(response); err != nil {
		fmt.Println("Error during json encoding in CreateCategory")
		return false
	} else {
		fmt.Println("Json sent successfully in CreateCategory")
		return true
	}
}
