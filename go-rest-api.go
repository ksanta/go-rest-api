package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ksanta/go-rest-api/domain"
	"github.com/ksanta/go-rest-api/repository"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var repo repository.EventRepository

func homeLink(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome! /events resource is available")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent domain.Event
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

	err = repo.Create(&newEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventId := mux.Vars(r)["id"]

	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := repo.GetById(eventIdInt)
	if event == nil && err == nil {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall event and write to response
	json.NewEncoder(w).Encode(event)
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall event and write to response
	json.NewEncoder(w).Encode(*events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var updatedEvent domain.Event
	json.Unmarshal(reqBody, &updatedEvent)
	updatedEvent.ID, _ = strconv.Atoi(eventID)

	err = repo.Update(&updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedEvent)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	eventIdInt, _ := strconv.Atoi(eventID)

	err := repo.Delete(eventIdInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	repo, err = repository.NewPostgresEventRepo()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
