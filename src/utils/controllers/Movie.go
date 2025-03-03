package controllers

import (
	"encoding/json"
	"io"
	"main/src/utils/constants"
	"main/src/utils/dao"
	"main/src/utils/models"
	"net/http"

	"github.com/gorilla/mux"
)

// Expects Theatre details so provide all movies available in theatres
func GetMovieList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["theatreId"]

	rows, err := dao.GetDbObject().Query(constants.SPECIFIC_THEATRE_MOVIES, id)

	if err != nil {
		http.Error(w, "Failed Loading DB Movies", http.StatusInternalServerError)
		return
	}
	var AllMovies []models.MovieInfo

	for rows.Next() {
		var movieObj models.MovieInfo
		err := rows.Scan(&movieObj.Id, &movieObj.MovieName, &movieObj.MovieDisc, &movieObj.MovieRating, &movieObj.Time, &movieObj.TheaterId)
		if err != nil {
			http.Error(w, "Failed Loading DB Movies", http.StatusInternalServerError)
			return
		}

		AllMovies = append(AllMovies, movieObj)
	}

	response := map[string]any{
		"Status":       "Success",
		"AllMovieInfo": AllMovies,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	var MovieInfo models.MovieInfo
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed Reading BODY", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &MovieInfo)

	res, err := dao.GetDbObject().Exec(constants.INSERT_MOVIE, MovieInfo.MovieName, MovieInfo.MovieDisc, MovieInfo.MovieRating, MovieInfo.Time, MovieInfo.TheaterId)
	if err != nil {
		http.Error(w, "Failed Adding Movie", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"Status":       "Success",
		"DBResult": res,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}
