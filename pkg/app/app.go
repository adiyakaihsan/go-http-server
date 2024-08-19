package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/config"
)

type App struct {
	db *sql.DB
}

func connectDB() (*sql.DB, error)  {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
				config.Db_host, config.Db_port, config.Db_username, config.Db_password, config.Db_name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	log.Printf("Successfully connected to Database")
	return db, err
}

func Run() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Error Opening database: %v", err)
	}
	defer db.Close()

    fmt.Println("Hello, Go!")
	fmt.Println("Starting server on port 8080")
	
	app := App{}
	app.db = db

	// TODO: https://github.com/julienschmidt/httprouter
	// TODO: REST API
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/createUser",app.createUserHandler)
	http.HandleFunc("/getUser", app.getUserHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server")
	}
}
