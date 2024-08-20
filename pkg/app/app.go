package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/config"
	"github.com/julienschmidt/httprouter"
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

	app := App{}
	app.db = db

	router := httprouter.New()

	router.GET("/", rootHandler)
	router.GET("/getUser", app.getUserHandler)
	router.POST("/createUser", app.createUserHandler)

	fmt.Println("Starting server on port 8080")
	
	// TODO: REST API

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server")
	}
}
