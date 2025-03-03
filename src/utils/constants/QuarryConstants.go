package constants

const (

	// 	id int auto_increment primary key,
	//     username varchar(255) not null,
	//     email varchar(255) not null,
	//     password varchar(500) not null

	INSERT_USER = "INSERT INTO users(username,email,password) VALUES (?,?,?);"
	REMOVE_USER = "DELETE FROM users where email = ? AND password = ?"
	VERIFY_USER = "SELECT username,email from users where email = ? AND PASSWORD = ?"

	// 	id int auto_increment primARY KEY,
	// 	TheaterName varchar(255) not null,
	// 	Location varchar(255) not null,
	// 	TheatreDisc varchar(255) not null

	INSERT_THEATRE = "INSERT INTO Theaters(TheaterName,Location,TheatreDisc) VALUES (?,?,?);"
	REMOVE_THEATRE = "DELETE FROM Theaters where id = ?;"
	ALL_THEATRE = "SELECT * FROM Theaters;"

	// 	id INT AUTO_INCREMENT PRIMARY KEY,
	// 	MovieName VARCHAR(255) NOT NULL,
	// 	MovieDisc VARCHAR(255) NOT NULL,
	// 	MovieRating INT NOT NULL,
	// 	Time DATETIME NOT NULL,
	// 	TheaterId INT NOT NULL

	INSERT_MOVIE = "INSERT INTO MovieInfo(MovieName,MovieDisc,MovieRating,Time,TheaterId) VALUES (?,?,?,?,?);"
	REMOVE_MOVIE = "DELETE FROM MovieInfo where id = ?"

	// 	    id INT AUTO_INCREMENT PRIMARY KEY,
	// 		Letter varchar(255) not null,
	//         SeatNum int not null,
	//         Side varchar(255) not null,
	//         MovieId int not null

	INSERT_SEAT = "INSERT INTO Seat(Letter,SeatNum,Side,MovieId) VALUES (?,?,?,?);"
	REMOVE_SEAT = "DELETE FROM Seat where id = ?"
)
