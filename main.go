package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	db_host = "localhost"
	db_port = 5432
	db_name = "production"
	db_username = "app_go"
	db_password = "app12345"
)

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
	}
	w.Write(body)
}

func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Printf("%s: %v", "Invalid JSON", err)
			return
		}

		sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
		
		_, err = db.Exec(sqlStatement, user.Username, user.Password)

		if err != nil {
			fmt.Print("here")
			http.Error(w, "Error inserting user", http.StatusInternalServerError)
			log.Printf("%s: %v", "Error inserting user", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
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
		// QUESTION: why some method in GO doesn't have err return like below?
		row := db.QueryRow(sqlStatement, user["username"])

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
	
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
				db_host, db_port, db_username, db_password, db_name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error Opening database: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to Database")

    fmt.Println("Hello, Go!")
	fmt.Println("Starting server on port 8080")

	http.HandleFunc("/", handler)
	http.HandleFunc("/createUser",createUserHandler(db))
	http.HandleFunc("/getUser", getUserHandler(db))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server")
	}
}