package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initCluster()

	r := mux.NewRouter()
	r.HandleFunc("/api/books", createBookHandler).Methods("POST")
	r.HandleFunc("/api/authors", createAuthorHandler).Methods("POST")
	r.HandleFunc("/api/books", getBooksHandler).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBookHandler).Methods("GET")
	r.HandleFunc("/api/books/{id}", deleteBookHandler).Methods("DELETE")
	r.HandleFunc("/api/books/{id}", updateBookHandler).Methods("PUT")
	if isMaster {
		go log.Fatal(http.ListenAndServe(":9000", r))
	}
}
