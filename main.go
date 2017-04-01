package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

//RSA Keys and Initialisation

const (
	privateKeyPath = "key/app.rsa"
	publicKeyPath  = "key/app.rsa.pub"
)

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func initKeys() {
	signKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key file")
		return
	}

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	verifyKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("error reading public key file")
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyBytes)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
}

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
	//Public Endpoints
	http.HandleFunc("/login", LoginHandler)

	//Protected Endpoints
	http.Handle("/resource/", negroni.New(negroni.HandlerFunc(ValidateTokenMiddleware), negroni.Wrap(http.HandlerFunc(ProtectedHandler))))

	http.Handle("/people", negroni.New(negroni.HandlerFunc(ValidateTokenMiddleware), negroni.Wrap(http.HandlerFunc(GetPeopleEndPoint))))

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

	http.ListenAndServe(":8000", nil)

}

func initData() {
	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Bangalore", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Lastname: "Patankar", Address: &Address{City: "San Fransico", State: "California"}})
	people = append(people, Person{ID: "3", Firstname: "Hun", Lastname: "Patankar", Address: &Address{City: "Munich", State: "Germany"}})
}
func main() {
	initKeys()
	initData()
	StartServer()
}

//EndPoint Handlers
func GetPeopleEndPoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	tokenString, err := token.SignedString(SignKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error Signing token %v\n", err)
	}

	//create a token instance using the token string
	response := Token{tokenString}
	JsonResponse(response, w)

}
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}
}

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
