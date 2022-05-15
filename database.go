package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book Struct
type Book struct {
	ID       int    `json:"ID" gorm:"autoIncrement" gorm:"primaryKey"`
	Isbn     string `json:"Isbn"`
	Title    string `json:"Title"`
	AuthorID int    `json:"AuthorID"`
}

// author struct

type Author struct {
	ID    int    `json:"ID" gorm:"primaryKey" gorm:"AutoIncrement"`
	First string `json:"First"`
	Last  string `json:"Last"`
}

type Event struct {
	ID   int    `json:"ID" gorm:"primaryKey" gorm:"AutoIncrement"`
	Code string `json:"Code"`
	Body string `json:"Body"`
}

type StatusMessage struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

var dbManager *gorm.DB

// connect to database
func connectDatabase(c chan string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot connect to database")
	}
	dbManager = db
	dbManager.AutoMigrate(&Book{}, &Author{}, &Event{})
	c <- "Connected to Database"
	return db
}

// author related functions
func createAuthor(author Author) (bool, Author) {
	result := dbManager.Create(&author)
	if result.Error == nil {
		return true, author
	}
	return false, author
}

// book related functions

func createBook(book Book) (bool, Book) {
	result := dbManager.Create(&book)
	if result.Error == nil {
		return true, book
	}
	return false, book
}

func deleteBook(id string) bool {
	book := Book{}
	err := dbManager.Delete(&book, id).Error
	if err != nil {
		return false
	}
	return true
}

func updateBook(book Book) (bool, Book) {
	fmt.Println(book)
	result := dbManager.Save(&book)
	fmt.Println(result)
	if result.Error == nil {
		return true, book
	}
	return false, book
}

func BookExist(id string) (bool, Book) {
	bookFound := Book{}
	err := dbManager.First(&bookFound, id).Error
	if err == nil {
		return true, bookFound
	} else {
		return false, bookFound
	}
}

func authorExist(id string) (bool, Author) {
	authorFound := Author{}
	err := dbManager.First(&authorFound, id).Error
	if err == nil {
		return true, authorFound
	} else {
		return false, authorFound
	}
}

func saveEvent(code string, body string) bool {
	result := dbManager.Create(&Event{Body: body, Code: code})
	if result.Error == nil {
		return true
	}
	return false
}
