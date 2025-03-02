-- All SQL Table Quarrys will be Listed Here


-- git config --global user.email "eswar.vallabha@iouring.com"

-- git config --global user.email "ramagownieswar@karunya.edu.in"



use MovieTicketBooking;

create table Theaters (
	id int auto_increment primARY KEY,
    TheaterName varchar(255) not null,
    Location varchar(255) not null,
    TheatreDisc varchar(255) not null
);


CREATE TABLE MovieInfo (
    id INT AUTO_INCREMENT PRIMARY KEY,
    MovieName VARCHAR(255) NOT NULL,
    MovieDisc VARCHAR(255) NOT NULL,
    MovieRating INT NOT NULL,
    Time DATETIME NOT NULL,
    TheaterId INT NOT NULL
);

CREATE TABLE Seat (
	    id INT AUTO_INCREMENT PRIMARY KEY,
		Letter varchar(255) not null,
        SeatNum int not null,
        Side varchar(255) not null,
        MovieId int not null 
);


show tables;

desc movieinfo;