package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashwinipatankar/JwtExInGo/authentication"
	handler "github.com/ashwinipatankar/JwtExInGo/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Struct Definitions

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Person struct {
	ID        string   `json:"id, omitempty"`
	Firstname string   `json:"firstname , omitempty"`
	Lastname  string   `json:"lastname, omitempty"`
	Address   *Address `json:"address, omitempty"`
}

type Address struct {
	City  string `json:"city, omitempty"`
	State string `json:"state, omitempty"`
}

var people []Person

//Server Entry Point
func StartServer() {
	r := mux.NewRouter()

	//Public Endpoints
	r.Handle("/", handler.GetLoginPageHandler).Methods("GET")
	r.Handle("/login", handler.LoginHandler).Methods("POST")

	//Protected Endpoints

	r.Handle("/people", authentication.ValidateToken.Handler(GetPeopleEndPointHandler)).Methods("GET")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(GetPersonEndPointHandler)).Methods("GET")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(CreatePersonEndPointHandler)).Methods("POST")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(DeletePersonEndPointHandler)).Methods("DELETE")

	log.Println("Now listening...")

	//handle server interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Stopping Server ...")
		os.Exit(1)
	}()

	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r))
}

func initData() {
	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Bangalore", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Lastname: "Patankar", Address: &Address{City: "San Fransico", State: "California"}})
	people = append(people, Person{ID: "3", Firstname: "Hun", Lastname: "Patankar", Address: &Address{City: "Munich", State: "Germany"}})
}

func main() {
	authentication.InitKeys()
	initData()
	StartServer()
}

//Endpoints
var GetPeopleEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
})

var GetPersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})

})

var CreatePersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)

})

var DeletePersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
})
