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

func main() {
	var err error
	repo, err = repository.NewPostgresEventRepo()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", welcomePage)
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

var repo repository.EventRepository

func welcomePage(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome! /events resource is available")
	if err != nil {
		log.Fatal(err)
	}
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent domain.Event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(reqBody, &newEvent)
	if err != nil {
		http.Error(w, "Please put an event object in the request body", http.StatusBadRequest)
		return
	}

	if err = newEvent.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = repo.Create(&newEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newEvent)
	if err != nil {
		log.Fatal(err)
	}
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
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllEvents(w http.ResponseWriter, _ *http.Request) {
	events, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall event and write to response
	err = json.NewEncoder(w).Encode(*events)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var updatedEvent domain.Event
	err = json.Unmarshal(reqBody, &updatedEvent)
	if err != nil {
		http.Error(w, "Please put an event object in the request body", http.StatusBadRequest)
		return
	}

	updatedEvent.ID, err = strconv.Atoi(eventID)
	if err != nil {
		http.Error(w, "ID is not a valid integer", http.StatusBadRequest)
	}

	// todo: validate to see if the data already exists

	err = repo.Update(&updatedEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(updatedEvent)
	if err != nil {
		log.Fatal(err)
	}
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
