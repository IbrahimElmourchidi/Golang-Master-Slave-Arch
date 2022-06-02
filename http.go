package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// create an Author
func createAuthorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var author Author
	_ = json.NewDecoder(r.Body).Decode(&author)
	authorStruct := Author{}
	ok := mapstructure.Decode(author, &authorStruct)
	if ok != nil {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
	dbOk, dbAuthor := createAuthor(authorStruct)
	if dbOk {
		message := tcpMessage{
			Code: "001",
			Body: serialize(dbAuthor),
		}
		saveEvent(message.Code, serialize(message.Body))
		go messageToAll(message)
		json.NewEncoder(w).Encode(dbAuthor)
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
}

// create a book
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	bookStruct := Book{}
	ok := mapstructure.Decode(book, &bookStruct)
	if ok != nil {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal server error",
		})
	}
	authorExist, foundAuthor := authorExist(strconv.Itoa(bookStruct.AuthorID))
	if authorExist {
		bookStruct.Author = foundAuthor
		dbOk, dbBook := createBook(bookStruct)
		if dbOk {
			message := tcpMessage{
				Code: "002",
				Body: serialize(dbBook),
			}
			saveEvent(message.Code, serialize(message.Body))
			go messageToAll(message)
			json.NewEncoder(w).Encode(dbBook)
		} else {
			json.NewEncoder(w).Encode(StatusMessage{
				Code:    "500",
				Message: "Internal Server Error",
			})
		}
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "404",
			Message: "Author Not Found",
		})
	}
}

// get all books
func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var allBooks []Book
	err := dbManager.Preload("Author").Find(&allBooks).Error
	if err == nil {
		json.NewEncoder(w).Encode(allBooks)
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
}

// get all authors
func getAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var allAuthors []Author
	err := dbManager.Find(&allAuthors).Error

	if err == nil {
		json.NewEncoder(w).Encode(allAuthors)
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
}

// get one books
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	bookExist, bookFound := BookExist(params["id"])
	if bookExist {
		json.NewEncoder(w).Encode(bookFound)
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "404",
			Message: "Book Not Found",
		})
	}
}

// delete book
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	exist, dbBook := BookExist(params["id"])
	if exist == false {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "404",
			Message: "Book Not Found",
		})
		return
	}
	idInt, _ := strconv.Atoi(params["id"])
	dbOk := deleteBook(idInt)
	if dbOk {
		message := tcpMessage{
			Code: "004",
			Body: serialize(dbBook),
		}
		saveEvent(message.Code, serialize(message.Body))
		go messageToAll(message)
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "202",
			Message: "Deleted Successfuly",
		})
		return
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
}

// update book
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	bookStruct := Book{}
	ok := mapstructure.Decode(book, &bookStruct)
	if ok != nil {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "500",
			Message: "Internal Server Error",
		})
		return
	}
	// check if book exsist & author exist
	bookExist, bookFound := BookExist(params["id"])
	authorExist, _ := authorExist(strconv.Itoa(bookStruct.AuthorID))
	if bookExist && authorExist {
		bookFound = assign(bookFound, bookStruct)
		dbOk, dbBook := updateBook(bookFound)
		if dbOk {
			message := tcpMessage{
				Code: "003",
				Body: serialize(dbBook),
			}
			saveEvent(message.Code, serialize(message.Body))
			go messageToAll(message)
			json.NewEncoder(w).Encode(dbBook)
			return
		} else {
			json.NewEncoder(w).Encode(StatusMessage{
				Code:    "500",
				Message: "Internal Server Error",
			})
			return
		}
	} else {
		json.NewEncoder(w).Encode(StatusMessage{
			Code:    "404",
			Message: "Book or Author not exist",
		})
		return
	}

}

func assign(original Book, final Book) Book {
	original.Isbn = final.Isbn
	original.Title = final.Title
	original.AuthorID = final.AuthorID
	original.Price = final.Price
	return original
}

func serialize(body interface{}) string {
	return fmt.Sprintf("%v", body)
}
