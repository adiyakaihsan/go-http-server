package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adiyakaihsan/go-http-server/pkg/types"
	"github.com/julienschmidt/httprouter"
)

func (app App) createVideoHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var video types.Video

	err := json.NewDecoder(r.Body).Decode(&video)

	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		log.Printf("%s: %v", "Unable to read request body", err)
		return
	}

	sqlStatement := `INSERT INTO videos (title, description, category_id) VALUES ($1, $2, $3)`

	_, err = app.db.Exec(sqlStatement, video.Title, video.Description, video.CategoryID)

	if err != nil {
		http.Error(w, "Error inserting video", http.StatusInternalServerError)
		log.Printf("%s: %v", "Error inserting video", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app App) getAllVideosHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var videosResponse []types.VideoResponse

	rows, err := app.db.Query("Select id, title, description, category_id FROM videos")
	if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var videoResponse types.VideoResponse
		err := rows.Scan(&videoResponse.ID, &videoResponse.Title, &videoResponse.Description, &videoResponse.CategoryID)

		if err != nil {
			http.Error(w, "Error Scanning Rows", http.StatusInternalServerError)
			log.Printf("Error Scanning Rows: %v", err)
			return
		}
		videosResponse = append(videosResponse, videoResponse)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		log.Printf("Error iterating rows: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(videosResponse)
}

func (app App) getVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var videoResponse types.VideoResponse

	id := ps.ByName("id")

	err := app.db.QueryRow("SELECT id, title, description, category_id  FROM videos WHERE id = $1", id).Scan(&videoResponse.ID, &videoResponse.Title, &videoResponse.Description, &videoResponse.CategoryID)

	if err == sql.ErrNoRows {
		http.Error(w, "video Not Found", http.StatusNotFound)
		log.Printf("video with ID: %v not found", id)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving Data", http.StatusInternalServerError)
		log.Printf("Error retrieving Data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(videoResponse)

}
