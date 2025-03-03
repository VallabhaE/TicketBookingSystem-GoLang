package gameLoader

import (
	"main/src/utils/controllers"
	"main/src/utils/dao"
	"net/http"

	"github.com/gorilla/mux"
)

func __init() {
	dao.InitDB()
	dao.Ping()
}

func Start() {
	__init()

	r := mux.NewRouter()

	// User Level Routes
	r.HandleFunc("/login",controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup",controllers.SignUp).Methods(http.MethodPost)


	// Theatre Level Route
	r.HandleFunc("/allTheatres",controllers.GetAllTheatres).Methods(http.MethodGet)
	r.HandleFunc("/addTheatre",controllers.AddTheatre).Methods(http.MethodPost)

	// Movie Level Route
	


	http.Handle("/",r)

	
}	