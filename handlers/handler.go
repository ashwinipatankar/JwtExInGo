package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ashwinipatankar/JwtExInGo/authentication"
	"github.com/ashwinipatankar/JwtExInGo/data"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//App claims provide custom claim for JWt
type AppClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}

var GetLoginPageHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/login.html")
})

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

var GetPeopleEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(data.GetPeople())
})

var GetPersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range data.GetPeople() {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&data.Person{})

})

var CreatePersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person data.Person
	people := data.GetPeople()
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(data.GetPeople())

})

var DeletePersonEndPointHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	people := data.GetPeople()
	defer data.SetPeople(people)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
})
