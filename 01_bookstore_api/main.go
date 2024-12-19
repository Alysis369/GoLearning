package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Create data model
type Book struct {
	ID     int    `json:"id"` // backticks is to define what the JSON title will be, ex. have to be in quotation and string
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Create slice of books
var books []Book

func preloadBooks() {
	books = append(books, Book{ID: 1, Title: "The Great Aldo", Author: "Aldo S"})
}

/* General API handler without mux */
func handleBooks(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)		// No longer needed when using json.Encode(), Go automatically set header status to OK if content is written
		// w.Write([]byte(`[]`)) 			// Use backticks for string literals \n won't be interpreted
		json.NewEncoder(w).Encode(books)

	case r.Method == http.MethodPost:
		var newBook Book
		json.NewDecoder(r.Body).Decode(&newBook)
		newBook.ID = len(books) + 1
		books = append(books, newBook)

		// Return to client declaring the newly created book
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBook)
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newBook.ID = len(books) + 1
	books = append(books, newBook)

	// Return to client declaring the newly created book
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Decode the requestedbook
	var updatedBook Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook) // REMEMBER DECODE TAKES POINTER
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Find the book
	for i, book := range books {
		if fmt.Sprintf("%d", book.ID) == id { // Use sprintf since in Go, string(<int>) will assume the int is the unicode code point
			// Update book
			books[i] = updatedBook

			// Send back confirmation that book is sent
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	// Reached when book is not found
	http.Error(w, "Book not found", http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Search books for id
	for i, book := range books {
		if fmt.Sprintf("%d", book.ID) == id {
			// delete the book
			books = append(books[:i], books[i+1:]...) // ... is variadic operator, unpacks the slice to individual elem
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Reached when book is not found
	http.Error(w, "Book not found", http.StatusNotFound)
}

func main() {
	fmt.Println("Server is running...")
	preloadBooks() // Preload books

	r := mux.NewRouter()

	// http.HandleFunc("/books", handleBooks)	// Single direction, change format to use mux
	r.HandleFunc("/books", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PATCH")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	http.ListenAndServe("localhost:8080", r) // Pass in the mux here, to tell to route through created mux
}
