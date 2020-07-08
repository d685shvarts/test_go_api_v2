package main

import (
	//Utility libraries
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//Router
	"github.com/gorilla/mux"

	//MySql
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

var notes []Note

func main() {
	notes = append(notes, Note{ID: "1", Title: "Test Note", Text: "Just a simple test note"})
	//This is the router
	router := mux.NewRouter()

	//Start MySQL database

	db, _ := sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/noteDatabase")

	if err := db.Ping(); err != nil {
		db.Close()
		fmt.Printf("Error pinging DB")
		return
	}
	defer db.Close()

	//Here we create all the endpoints
	//HandleFunc takes two arguments: a string that defiens the route and a function that handle the roue
	//.method() takes an argument of the http method
	//THis is all the methods we need for CRUD (Create, Read, Update, Delete)
	router.HandleFunc("/note", getNotes).Methods("GET")
	router.HandleFunc("/note/{id}", getNote).Methods("GET")
	router.HandleFunc("/note", createNote).Methods("POST")
	router.HandleFunc("/note/{id}", updateNote).Methods("POST")
	router.HandleFunc("/note/{id}", deleteNote).Methods("DELETE")

	//Here ListenAndServe actually runs the server, its first argument is the address of the server and
	//second argument is the handler for the address
	//We wrap it all in the log.Fatal() to throw an error if it fails
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	//This sets the header of the response
	w.Header().Set("Content-Type", "application/json")
	//Encode then encodes our notes slice as json and sends it to response stream
	//Encoders write JSOn to outpput stream
	//NEw Encoder returns a new encoder that writes to w
	json.NewEncoder(w).Encode(notes)

}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get variable from http repsonse we are passing to it
	params := mux.Vars(r)

	for _, item := range notes {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//First create a new Note struct to store the request
	var newNote Note
	//To get the data of the json response, we must Decoder (similar to using encoder to put into r)
	//NewDecoder returns a new decoder that reads from r
	//Decode reads the next JSON encoded value from input and store it in newNote
	json.NewDecoder(r.Body).Decode(&newNote)

	fmt.Printf("%+v\n", newNote)
	//Here we simulate an id by getting the length of the notes slice, then using strConv Itoa
	//to convert it to a string
	newNote.ID = strconv.Itoa(len(notes) + 1)

	//Add new roll to struct to our notes slice
	notes = append(notes, newNote)

	//Then send response back containing new roll
	json.NewEncoder(w).Encode(newNote)
}

func updateNote(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//Get variable from http repsonse we are passing to it
	params := mux.Vars(r)

	for i, item := range notes {
		if item.ID == params["id"] {
			//This removes the current item from the notes slice
			notes = append(notes[:i], notes[i+1:]...)
			var newNote Note
			json.NewDecoder(r.Body).Decode(&newNote)
			newNote.ID = params["id"]
			notes = append(notes, newNote)
			json.NewEncoder(w).Encode(newNote)
			return
		}
	}

}
func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get variable from http repsonse we are passing to it
	params := mux.Vars(r)

	for i, item := range notes {
		if item.ID == params["id"] {
			//This removes the current item from the notes slice
			notes = append(notes[:i], notes[i+1:]...)

			break
		}
	}

	json.NewEncoder(w).Encode(notes)

}
