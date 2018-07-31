package main

import (
		_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

type Toilet struct {
	gorm.Model
	Name string
	Location string
	ToiletRatings []ToiletRating
}

type ToiletRating struct {
	gorm.Model
	ToiletId int
	Toilet Toilet
	Rating int
	Comment string
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalln("Failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Toilet{}, &ToiletRating{})

	runServer()
}
