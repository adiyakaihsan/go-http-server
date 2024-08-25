package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/types"
	"github.com/julienschmidt/httprouter"
)

func (app App) createCategoryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var category types.Category

	err := json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	sqlStatement := `INSERT INTO categories (name) VALUES ($1)`

	_, err = app.db.Exec(sqlStatement, category.Name)

	if err != nil {
		http.Error(w, "Error inserting category", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error inserting category", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app App) getAllCategoriesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var categories []types.Category

	rows, err := app.db.Query("Select id, name FROM categories")
	if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var category types.Category
		err := rows.Scan(&category.ID, &category.Name)

		if err != nil {
			http.Error(w, "Error Scanning Rows", http.StatusInternalServerError)
			log.Printf("Error Scanning Rows: %v", err)
			return
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		log.Printf("Error iterating rows: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(categories)
}

func (app App) getCategoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var category types.Category

	id := ps.ByName("id")

	err := app.db.QueryRow("SELECT id, name  FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)

	if err == sql.ErrNoRows {
		http.Error(w, "category Not Found", http.StatusNotFound)
		log.Printf("category with ID: %v not found", id)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(category)

}
