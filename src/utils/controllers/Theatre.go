package controllers

import (
	"encoding/json"
	"io"
	"main/src/utils/constants"
	"main/src/utils/dao"
	"main/src/utils/models"
	"net/http"
)

func GetAllTheatres(w http.ResponseWriter, r *http.Request) {

	res, err := dao.GetDbObject().Query(constants.ALL_THEATRE)
	if err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
	var allTheatres []models.Theaters
	for res.Next() {
		var singleTheatre models.Theaters
		err := res.Scan(&singleTheatre.Id, &singleTheatre.TheaterName, &singleTheatre.Location, &singleTheatre.TheatreDisc)

		if err != nil {
			http.Error(w, "Failed Quarrying DB", http.StatusInternalServerError)
			return
		}

		allTheatres = append(allTheatres, singleTheatre)
	}

	response := map[string]any{
		"Status":  "Success",
		"AllInfo": allTheatres,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}


func AddTheatre(w http.ResponseWriter, r *http.Request){
	var Theatre models.Theaters
	data,err := io.ReadAll(r.Body)
	if err !=nil{
		http.Error(w, "Failed Reading BODY", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data,&Theatre)

	if err !=nil{
		http.Error(w, "Failed Extracting Body", http.StatusInternalServerError)
		return
	}

	res,err := dao.GetDbObject().Exec(constants.INSERT_THEATRE,Theatre.TheaterName,Theatre.Location,Theatre.TheatreDisc)

	if err !=nil{
		http.Error(w, "Failed To Add Theatre", http.StatusInternalServerError)
		return
	}
	response := map[string]any{
		"Status":  "Success",
		"DbInfo": res,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}