package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/types"
	"github.com/julienschmidt/httprouter"
)

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app App) createUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err = app.db.Exec(sqlStatement, user.Username, user.Password)

	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error inserting user", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app App) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user types.User

	id := ps.ByName("id")

	err := app.db.QueryRow("SELECT id, username FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)

	if err == sql.ErrNoRows {
		http.Error(w, "User Not Found", http.StatusNotFound)
		log.Printf("User with ID: %v not found", id)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)

}
