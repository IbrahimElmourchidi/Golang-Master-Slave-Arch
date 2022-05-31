package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	// Files are provided as a slice of strings.
	files := []string{
		"templates/home.html",
	}
	result := template.Must(template.New("home.html").ParseFiles(files...))
	errTemplate := template.Must(template.New("error.html").ParseFiles(files...))

	var allBooks []Book
	err := dbManager.Preload("Author").Find(&allBooks).Error

	if err == nil {
		result.Execute(w, allBooks)
	} else {
		errTemplate.Execute(w, "Data Base Error")
	}
}

func updateBookTemplate(w http.ResponseWriter, r *http.Request) {
	// Files are provided as a slice of strings.
	files := []string{
		"templates/update.html",
	}
	result := template.Must(template.New("update.html").ParseFiles(files...))
	errTemplate := template.Must(template.New("error.html").ParseFiles(files...))
	params := mux.Vars(r)
	bookExist, bookFound := BookExist(params["id"])
	if bookExist {
		result.Execute(w, bookFound)
	} else {
		errTemplate.Execute(w, "Book Not Found")
	}
}
