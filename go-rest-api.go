package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type event struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

var db *sql.DB

func homeLink(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome! /events resource is available")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &newEvent)
	if err != nil {
		http.Error(w, "Please put an event object in the request body", http.StatusBadRequest)
		return
	}

	if newEvent.Title == "" || newEvent.Description == "" {
		http.Error(w, "Title and description must be populated", http.StatusBadRequest)
	}

	result, err := db.Exec("insert into events (title, description) VALUES ($1, $2)", newEvent.Title, newEvent.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastInsertId, _ := result.LastInsertId()
	newEvent.ID = int(lastInsertId)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	// Query from database
	row := db.QueryRow("SELECT id, title, description FROM events where id = $1", eventID)
	event := event{}
	err := row.Scan(&event.ID, &event.Title, &event.Description)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Fatal(err)
	}

	// Marshall event and write to response
	json.NewEncoder(w).Encode(event)
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	// Query from database
	rows, err := db.Query("SELECT id, title, description FROM events")
	if err != nil {
		log.Fatal(err)
	}
	// Load results into an events slice
	events := make([]event, 0)
	for rows.Next() {
		event := event{}
		err = rows.Scan(&event.ID, &event.Title, &event.Description)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}

	// Marshall event and write to response
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var updatedEvent event
	json.Unmarshal(reqBody, &updatedEvent)
	updatedEvent.ID, _ = strconv.Atoi(eventID)

	_, err = db.Exec("update events set title = $2, description = $3 where id = $1", updatedEvent.ID, updatedEvent.Title, updatedEvent.Description)
	if err != nil {
		w.Write([]byte("Error updating an existing event"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedEvent)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	_, err := db.Exec("delete from events where id = $1", eventID)
	if err != nil {
		w.Write([]byte("Error deleting an existing event"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func initialiseDatabaseConnection() {
	connStr := "user=postgres password=password dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initialiseDatabaseConnection()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
