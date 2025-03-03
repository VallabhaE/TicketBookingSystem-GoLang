package gameLoader

import (
	"fmt"
	"main/src/utils/controllers"
	"main/src/utils/dao"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	__init()

	r := mux.NewRouter()

	// User Level Routes
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", controllers.SignUp).Methods(http.MethodPost)

	// Theatre Level Route
	r.HandleFunc("/allTheatres", controllers.GetAllTheatres).Methods(http.MethodGet)
	r.HandleFunc("/addTheatre", controllers.AddTheatre).Methods(http.MethodPost)

	// Movie Level Route
	r.HandleFunc("/GetMovieList/{theatreId}", controllers.GetMovieList).Methods(http.MethodGet)
	r.HandleFunc("/AddMovie/", controllers.AddMovie).Methods(http.MethodPost)

	//Seats Level Routes
	// Steps:
	// 1.Get all seats available for that Movie in that perticular theatre
	// 2.Book ticket Func
	// 		-- This is Main And impotent Func we can say
	//		-- it should lock seat for 10 min and book seat if user pays money
	// 		-- Temp use Local Cache later add Redis for distributed System Ticket booking
	//		-- Add RazerPay Payment Integration



	r.HandleFunc("/GetAllSeats/{movieId}", controllers.GetAllSeats).Methods(http.MethodGet)
	r.HandleFunc("/LockSeat", controllers.LockSeat).Methods(http.MethodGet)
	r.HandleFunc("/VerifyPayment", controllers.VerifyPayment).Methods(http.MethodPost)

	fmt.Println("Payment Integration is Pending\nTicket Booking Logic using transactions\nremove local cache and add redis for centrilised system if needed")

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server: ", err)
	}

}

func __initDao() {
	dao.InitDB()
	dao.Ping()
}
func __init() {
	__initDao()
}
