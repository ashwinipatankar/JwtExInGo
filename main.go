package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ashwinipatankar/JwtExInGo/authentication"
	handler "github.com/ashwinipatankar/JwtExInGo/handlers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Struct Definitions
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
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

//App claims provide custom claim for JWt
type AppClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

//Server Entry Point
func StartServer() {
	r := mux.NewRouter()

	//Public Endpoints
	r.Handle("/", handler.GetLoginPageHandler).Methods("GET")
	r.Handle("/login", LoginHandler).Methods("POST")

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

/*
var GetLoginPageHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
})
*/
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	//decode request into user credentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	fmt.Println(user.Username, user.Password)

	//Validate user credentials
	if strings.Compare(strings.ToLower(user.Username), "admin") != 0 {
		if strings.Compare(user.Password, "password") != 0 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	//Create claims
	claims := AppClaims{user.Username, "Member", jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		Issuer:    "admin",
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(authentication.GetSignKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error Signing token %v\n", err)
	}

	//create a token instance using the token string
	response := Token{tokenString}
	JsonResponse(response, w)

})

//Helper Function
func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
