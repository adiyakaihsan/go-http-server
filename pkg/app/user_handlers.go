package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adiyakaihsan/go-http-server/pkg/config"
	"github.com/adiyakaihsan/go-http-server/pkg/types"
	jwt "github.com/golang-jwt/jwt/v4"
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

func generateJWTToken(id int, username string) []byte {
	// var JWT_SIGNING_METHOD jwt.SigningMethodHS256

	claims := types.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.APP_NAME,
			ExpiresAt: time.Now().Add(config.LOGIN_EXPIRATION_DURATION).Unix(),
		},
		ID:       id,
		Username: username,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	log.Printf("%s: %v", "Token", token)
	signedToken, err := token.SignedString(config.JWT_SIGNATURE_KEY)
	if err != nil {
		log.Printf("%s: %v", "Error when generating token", err)
		return nil
	}

	tokenResponse := types.TokenResponse{Token: signedToken}

	token_string, _ := json.Marshal(tokenResponse)

	return token_string
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

func (app App) loginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user types.User
	var hashedPassword string
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error when hashing password user", err)
		return
	}

	err = app.db.QueryRow("SELECT id,password FROM users WHERE username = $1", user.Username).Scan(&user.ID, &hashedPassword)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid Login", http.StatusNotFound)
		log.Printf("User with Username: %v invalid login pass: %v", user.Username, hashedPassword)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))

	if err != nil {
		http.Error(w, "Invalid Login", http.StatusNotFound)
		log.Printf("User with Username: %v invalid login pass: %v", user.Username, hashedPassword)
		return
	}

	jwt_token := generateJWTToken(user.ID, user.Username)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(jwt_token)
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
