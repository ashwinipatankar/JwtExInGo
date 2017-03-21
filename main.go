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
	params := mux.Vars(request)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(writer).Encode(item)
			return
		}
	}
	json.NewEncoder(writer).Encode(&Person{})

}

func GetPeopleEndPoint(writer http.ResponseWriter, request *http.Request) {
	json.NewEncoder(writer).Encode(people)

}

func CreatePersonEndPoint(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(writer).Encode(people)

}

func DeletePersonEndPoint(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(writer).Encode(people)
}
func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Hyderanad", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Address: &Address{City: "Bangalore", State: "India"}})
	people = append(people, Person{ID: "3", Firstname: "John"})
	router.HandleFunc("/people", GetPeopleEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndPoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndPoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
