package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/src/utils/constants"
	"main/src/utils/dao"
	"main/src/utils/models"
	"main/src/utils/razorpay"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllSeats(w http.ResponseWriter, r *http.Request) {
	var seats []models.Seat
	id := mux.Vars(r)["movieId"]

	rows, err := dao.GetDbObject().Query(constants.ALL_SEATS, id)
	if err != nil {
		http.Error(w, "Failed Loading DB Movies", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(&seat.Id, &seat.Letter, &seat.SeatNum, &seat.Side, &seat.MovieId)
		if err != nil {
			http.Error(w, "Failed Loading DB Seats to Local", http.StatusInternalServerError)
			return
		}
		seats = append(seats, seat)
	}

	response := map[string]any{
		"Status":   "Success",
		"AllSeats": seats,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}

}

func LockSeat(w http.ResponseWriter, r *http.Request) {
	// Extract seat_letter and seat_num from the request URL or form
	seatLetter := r.URL.Query().Get("letter")
	seatNum := r.URL.Query().Get("seat_num")

	// Validate input
	if seatLetter == "" || seatNum == "" {
		http.Error(w, "Both seat letter and seat number are required", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := dao.GetDbObject().BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Use a defer with explicit error check for the rollback
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but don't override the original error
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Step 1: Check if the seat exists and is already locked using QueryRow
	var locked bool
	err = tx.QueryRow(`SELECT locked FROM Seat WHERE Letter = $1 AND SeatNum = $2 FOR UPDATE`, seatLetter, seatNum).Scan(&locked)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Seat %s%s does not exist", seatLetter, seatNum), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to check seat availability: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Step 2: If the seat is already locked, return an error
	if locked {
		http.Error(w, fmt.Sprintf("Seat %s%s is already locked or taken by another user", seatLetter, seatNum), http.StatusConflict)
		return
	}

	// Step 3: Lock the seat by updating the "locked" column
	result, err := tx.Exec(`UPDATE Seat SET locked = TRUE WHERE Letter = $1 AND SeatNum = $2 AND locked = FALSE`, seatLetter, seatNum)
	if err != nil {
		http.Error(w, "Failed to lock seat: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get rows affected: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, fmt.Sprintf("Seat %s%s could not be locked - it may have been locked by another concurrent request", seatLetter, seatNum), http.StatusConflict)
		return
	}

	data, err := razorpay.CreateOrderId(100, int(rowsAffected))
	// Step 4: Commit the transaction if everything went well
	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": fmt.Sprintf("Seat %s%s has been successfully locked", seatLetter, seatNum),
		"data":    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
func VerifyPayment(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get payment verification details
	var paymentData struct {
		RazorpayOrderID    string `json:"razorpay_order_id"`
		RazorpayPaymentID  string `json:"razorpay_payment_id"`
		RazorpaySignature  string `json:"razorpay_signature"`
		SeatLetter         string `json:"seat_letter"`
		SeatNum            string `json:"seat_num"`
		UserID             int    `json:"user_id"` // Assuming you need to know which user made the booking
	}

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&paymentData); err != nil {
		http.Error(w, "Invalid request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if paymentData.RazorpayOrderID == "" || paymentData.RazorpayPaymentID == "" || 
		paymentData.RazorpaySignature == "" || paymentData.SeatLetter == "" || 
		paymentData.SeatNum == "" || paymentData.UserID == 0 {
		http.Error(w, "Missing required payment verification details", http.StatusBadRequest)
		return
	}

	// Verify the payment with Razorpay
	isPaymentValid := razorpay.VerifyPayment(
		paymentData.RazorpayOrderID, 
		paymentData.RazorpayPaymentID, 
		paymentData.RazorpaySignature,
	)
	
	// Start a transaction
	tx, err := dao.GetDbObject().BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use a defer with explicit error check for the rollback
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but don't override the original error
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Check if the seat is still locked
	var locked bool
	err = tx.QueryRow(`SELECT locked FROM Seat WHERE Letter = $1 AND SeatNum = $2 FOR UPDATE`, 
		paymentData.SeatLetter, paymentData.SeatNum).Scan(&locked)
	
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Seat %s%s does not exist", 
				paymentData.SeatLetter, paymentData.SeatNum), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to check seat status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !locked {
		http.Error(w, fmt.Sprintf("Seat %s%s is not locked. Booking session may have expired.", 
			paymentData.SeatLetter, paymentData.SeatNum), http.StatusConflict)
		return
	}

	// If payment is valid, confirm the booking
	if isPaymentValid {
		// // Insert booking record
		// _, err = tx.Exec(`
		// 	INSERT INTO Booking (UserID, SeatLetter, SeatNum, PaymentID, OrderID, BookingTime)
		// 	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		// `, paymentData.UserID, paymentData.SeatLetter, paymentData.SeatNum, 
		//    paymentData.RazorpayPaymentID, paymentData.RazorpayOrderID)
		
		// if err != nil {
		// 	http.Error(w, "Failed to create booking record: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		
		// Update seat status to booked (not just locked)
		_, err = tx.Exec(`
			UPDATE Seat 
			SET locked = TRUE, booked = TRUE, UserID = $1 
			WHERE Letter = $2 AND SeatNum = $3
		`, paymentData.UserID, paymentData.SeatLetter, paymentData.SeatNum)
		
		if err != nil {
			http.Error(w, "Failed to update seat status: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// If payment verification failed, unlock the seat
		_, err = tx.Exec(`
			UPDATE Seat 
			SET locked = FALSE, booked = FALSE, UserID = NULL 
			WHERE Letter = $1 AND SeatNum = $2
		`, paymentData.SeatLetter, paymentData.SeatNum)
		
		if err != nil {
			http.Error(w, "Failed to unlock seat: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return appropriate response based on payment status
	w.Header().Set("Content-Type", "application/json")
	
	if isPaymentValid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"message": fmt.Sprintf("Payment verified and seat %s%s has been successfully booked", 
				paymentData.SeatLetter, paymentData.SeatNum),
			"booking_details": map[string]interface{}{
				"seat": paymentData.SeatLetter + paymentData.SeatNum,
				"payment_id": paymentData.RazorpayPaymentID,
				"order_id": paymentData.RazorpayOrderID,
				"user_id": paymentData.UserID,
			},
		})
	} else {
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "failed",
			"message": "Payment verification failed. The seat lock has been released.",
		})
	}
}