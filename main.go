package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// fmt.Fprintf(w, "Hello World!")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body")
	}
	w.Write(body)
}

func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello World!")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Fatalf(err.Error())
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			log.Fatalf(err.Error())
			return
		}
		defer r.Body.Close()

		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Fatalf(err.Error())
			return
		}

		sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
		
		_, err = db.Exec(sqlStatement, user.Username, user.Password)

		if err != nil {
			http.Error(w, "Error inserting user", http.StatusInternalServerError)
			log.Fatalf(err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created Successfuly"))
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server")
	}
}