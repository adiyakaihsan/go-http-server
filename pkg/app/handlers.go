package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/types"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
	}
	w.Write(body)
}

func (app *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}
	defer r.Body.Close()

	var user types.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("%s: %v", "Invalid JSON", err)
		return
	}

	sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err = app.db.Exec(sqlStatement, user.Username, user.Password)

	if err != nil {
		fmt.Print("here")
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error inserting user", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) getUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Error reading request body", err)
	}

	defer r.Body.Close()
	var user map[string]interface{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("%s: %v", "Invalid JSON", err)
		return
	}

	sqlStatement := `SELECT username FROM users WHERE username=$1`
	row := app.db.QueryRow(sqlStatement, user["username"])

	var username string
	err = row.Scan(&username)
	if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(username)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Printf("Error marshaling response: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
