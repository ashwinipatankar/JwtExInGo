package main

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
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

var signingKey, verificationKey []byte

func intiKeys() {
	var (
		err         error
		privKey     *rsa.PrivateKey
		pubKey      *rsa.PublicKey
		pubKeyBytes []byte
	)

	privKey, err = rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		log.Fatal("Error generating private key")
	}
	pubKey = &privKey.PublicKey

	//Create signingkey from priv key
	//Prepare PEM block
	var privPEMBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey), //Marshal means searlize
	}

	//serialize pem
	privKeyPEMBuffer := new(bytes.Buffer)
	pem.Encode(privKeyPEMBuffer, privPEMBlock)

	//done
	signingKey = privKeyPEMBuffer.Bytes()

	fmt.Println(string(signingKey))

	//create verificationKey from pubKey, Also in PEM-format
	pubKeyBytes, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatal("Error marshalling public key")
	}

	var pubPEMBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	pubKeyPEMBuffer := new(bytes.Buffer)
	pem.Encode(pubKeyPEMBuffer, pubPEMBlock)

	//done
	verificationKey = pubKeyPEMBuffer.Bytes()

	fmt.Println(string(verificationKey))

}

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

func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user userCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Unauthorized access request / Error in Request")
		return
	}
	fmt.Println(user.Username, user.Password)

	//Integrate with Database
	if user.Username != "admin" || user.Password != "password" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "username/password doesnt match")
		fmt.Println("Error logging in because of username/password")
		return
	}

	//Create a rsa 256 signer
	//TODO: Add support to other ways by configuration file
	//signer := jwt.New(jwt.SigningMethodRS256)

	//set claims
	//things are broken by this, time to change the library

}
func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", GetLoginPage).Methods("GET")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndPoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndPoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndPoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}

func initData() {

	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Hyderanad", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Address: &Address{City: "Bangalore", State: "India"}})
	people = append(people, Person{ID: "3", Firstname: "John"})
}
func main() {

	initData()
	startServer()
}
