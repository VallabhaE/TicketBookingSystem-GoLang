package controllers

import (
	"encoding/json"
	"io"
	"main/src/utils/constants"
	"main/src/utils/dao"
	"main/src/utils/middleware"
	"main/src/utils/models"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var user models.User
	json.Unmarshal(data, &user)

	row := dao.GetDbObject().QueryRow(constants.VERIFY_USER, user.Email, user.Password)
	var verifiedUser models.User
	err = row.Scan(verifiedUser.Username, verifiedUser.Email)
	if err != nil {
		http.Error(w, "Failed to Verify User", http.StatusBadRequest)
		return
	}

	if verifiedUser.Username == "" || verifiedUser.Email == "" {
		http.Error(w, "User Not Exist", http.StatusBadRequest)
		return
	}
	response := map[string]string{
		"AuthCode": middleware.Jwt_sign(verifiedUser.Username),
		"username": verifiedUser.Username,
		"email":    verifiedUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	var user models.User
	json.Unmarshal(data, &user)

	if user.Password == "" {
		http.Error(w, "Password is Empty", http.StatusBadRequest)
		return
	}

	res, err := dao.GetDbObject().Exec(constants.INSERT_USER, user.Username, user.Email, user.Password)

	if err != nil {
		http.Error(w, "Failed To Add User"+err.Error(), http.StatusBadRequest)
		return
	}

	rest, _ := res.LastInsertId()
	response := map[string]string{
		"Status":       "Success",
		"LastInsertId": string(rest),
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}
