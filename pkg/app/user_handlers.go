package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/types"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (app App) createUserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error when hashing password user", err)
		return
	}

	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err = app.db.Exec(sqlStatement, user.Username, hashedPassword)

	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error inserting user", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app App) getAllUsersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users []types.User

	rows, err := app.db.Query("Select id, username FROM users")
	if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(&user.ID, &user.Username)

		if err != nil {
			http.Error(w, "Error Scanning Rows", http.StatusInternalServerError)
			log.Printf("Error Scanning Rows: %v", err)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		log.Printf("Error iterating rows: %v", err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
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
