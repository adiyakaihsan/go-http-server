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

func connectDB() (*sql.DB, error) {
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
	router.GET("/v1/users", app.getAllUsersHandler)
	router.GET("/v1/users/:id", app.getUserHandler)
	router.POST("/v1/users", app.createUserHandler)
	router.POST("/v1/users/login", app.loginUser)
	router.GET("/v1/videos", app.getAllVideosHandler)
	router.GET("/v1/videos/:id", app.getVideoHandler)
	router.POST("/v1/videos", app.createVideoHandler)
	router.GET("/v1/categories", app.getAllCategoriesHandler)
	router.GET("/v1/categories/:id", app.getCategoryHandler)
	router.POST("/v1/categories", app.createCategoryHandler)

	loggedRouter := middleware(router)

	fmt.Println("Starting server on port 8080")

	if err := http.ListenAndServe(":8080", loggedRouter); err != nil {
		fmt.Println("Error starting server")
	}
}
