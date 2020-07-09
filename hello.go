package main

import (
	//Utility libraries
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

var db *sql.DB
var err error

func main() {
	notes = append(notes, Note{ID: "1", Title: "Test Note", Text: "Just a simple test note"})

	//Start MySQL database
	db, _ = sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/noteDatabase")

	if err := db.Ping(); err != nil {
		db.Close()
		fmt.Printf("Error pinging DB")
		return
	}

	defer db.Close()

	//This is the router
	router := mux.NewRouter()

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
	w.Header().Set("Content-Type", "application/json")

	result, err := db.Query("SELECT notes_id, notes_title, notes_text FROM notes_tbl")
	if err != nil {
		fmt.Printf("Error making query")
		panic(err.Error())
	}
	defer result.Close()
	var get_notes []Note
	for result.Next() {
		var note Note
		err := result.Scan(&note.ID, &note.Title, &note.Text)
		if err != nil {
			fmt.Printf("Error scanning")
			panic(err.Error())
		}
		get_notes = append(get_notes, note)
	}
	json.NewEncoder(w).Encode(get_notes)

}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var note Note
	result, err := db.Query("SELECT notes_id, notes_title, notes_text FROM notes_tbl WHERE notes_id = ?", params["id"])
	if err != nil {
		fmt.Printf("Error making query")
		panic(err.Error())
	}

	for result.Next() {
		err := result.Scan(&note.ID, &note.Title, &note.Text)
		if err != nil {
			panic(err.Error())
		}
	}

	defer result.Close()
	json.NewEncoder(w).Encode(note)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var note Note

	json.NewDecoder(r.Body).Decode(&note)

	insForm, err := db.Prepare("INSERT INTO notes_tbl (notes_title, notes_text) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = insForm.Exec(note.Title, note.Text)

	result, err := db.Query("SELECT LAST_INSERT_ID()")

	if err != nil {
		panic(err.Error())
	}

	var last_id string
	for result.Next() {
		err := result.Scan(&last_id)
		if err != nil {
			panic(err.Error())
		}
	}

	note.ID = last_id

	json.NewEncoder(w).Encode(note)
}

func updateNote(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var note Note
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&note)
	fmt.Printf("%+v\n", note)

	stmt, err := db.Prepare("UPDATE notes_tbl SET notes_text = ? WHERE notes_id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(note.Text, params["id"])
	if err != nil {
		panic(err.Error())
	}

	//TODO: Addd some error checking to return error if ID doesn't exist in DB
	result, err := db.Query("SELECT notes_id, notes_title, notes_text FROM notes_tbl WHERE notes_id = ?", params["id"])
	if err != nil {
		fmt.Printf("Error making query")
		panic(err.Error())
	}

	for result.Next() {
		err := result.Scan(&note.ID, &note.Title, &note.Text)
		if err != nil {
			panic(err.Error())
		}
	}

	defer result.Close()
	json.NewEncoder(w).Encode(note)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM notes_tbl WHERE notes_id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])

}
