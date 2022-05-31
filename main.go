package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initCluster()
	// ajax
	r := mux.NewRouter()
	r.HandleFunc("/api/books", createBookHandler).Methods("POST")
	r.HandleFunc("/api/authors", createAuthorHandler).Methods("POST")
	r.HandleFunc("/api/books", getBooksHandler).Methods("GET")
	r.HandleFunc("/api/authors", getAuthorsHandler).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBookHandler).Methods("GET")
	r.HandleFunc("/api/books/{id}", deleteBookHandler).Methods("DELETE")
	r.HandleFunc("/api/books/{id}", updateBookHandler).Methods("PUT")

	// template
	r.HandleFunc("/home", HomePage).Methods("GET")
	r.HandleFunc("/update/{id}", updateBookTemplate).Methods("GET")

	fs := http.FileServer(http.Dir("./templates/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	if isMaster {
		go log.Fatal(http.ListenAndServe(":9000", r))
	}
}
