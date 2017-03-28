package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codegangsta/negroni"
)

//RSA Keys and Initialisation

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte

func initKeys() {
	var err error

	SignKey, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	VerifyKey, err = ioutil.ReadFile(publicKeyPath)
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
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//Server Entry Point
func StartServer() {
	//Public Endpoints
	http.HandleFunc("/login", LoginHandler)

	//Protected Endpoints
	http.Handle("/resource/", negroni.New(negroni.HandlerFunc(ValidateTokenMiddleware), negroni.Wrap(http.HandlerFunc(ProtectedHandler))))

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
func main() {
	initKeys()
	StartServer()
}

//EndPoint Handlers

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
	if strings.ToLower(user.Username) != "admin" {
		if user.Password != "password" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}
}
