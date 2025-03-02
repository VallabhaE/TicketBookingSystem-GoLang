package models

import "time"

// Theater Model represents each theater with its details
type Theaters struct {
	Id          int    `json:"id"`
	TheaterName string `json:"TheaterName"`
	Location    string `json:"Location"`
	TheatreDisc string `json:"TheatreDisc"`
}

// MovieInfo represents details of an individual movie being shown at a theater
type MovieInfo struct {
	Id          int       `json:"id"`
	MovieName   string    `json:"MovieName"`
	MovieDisc   string    `json:"MovieDisc"`
	MovieRating int       `json:"MovieRating"`
	Time        time.Time `json:"Time"`
	TheaterId int  `json:"TheaterId"`

}

// Seats model represents a list of available seats for a specific movie
type Seats struct {
	Id        int  `json:"id"`
	TheaterId int  `json:"TheaterId"`
	MovieId   int  `json:"MovieId"`
	Seats     []Seat `json:"Seat"`
}

// Seat model represents details of a specific seat in a theater
type Seat struct {
	Id        int  `json:"id"`
	SeatNum int    `json:"SeatNum"`
	MovieId int `json:"MovieId"`
	Side    string `json:"Side"`
	Letter  string `json:"SeatLetter"`
}
