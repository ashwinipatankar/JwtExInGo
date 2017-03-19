package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id, omitempty"`
	Firstname string   `json:"firstname, omitempty"`
	Lastname  string   `json:"lastname, omitempty"`
	Address   *Address `json:"address, omitempty"`
}

type Address struct {
	City  string `json:"city, omitempty"`
	State string `json:"state, omitempty"`
}

var people []Person

func GetPersonEndPoint(writer http.ResponseWriter, request *http.Request) {

}

func GetPeopleEndPoint(writer http.ResponseWriter, request *http.Request) {
	json.NewEncoder(writer).Encode(people)

}

func CreatePersonEndPoint(writer http.ResponseWriter, request *http.Request) {

}

func DeletePersonEndPoint(writer http.ResponseWriter, request *http.Request) {

}
func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Hyderanad", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Address: &Address{City: "Bangalore", State: "India"}})
	router.HandleFunc("/people", GetPeopleEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndPoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndPoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
