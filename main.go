package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashwinipatankar/JwtExInGo/authentication"
	"github.com/ashwinipatankar/JwtExInGo/data"
	handler "github.com/ashwinipatankar/JwtExInGo/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Server Entry Point
func StartServer() {
	r := mux.NewRouter()

	//Public Endpoints
	r.Handle("/", handler.GetLoginPageHandler).Methods("GET")
	r.Handle("/login", handler.LoginHandler).Methods("POST")

	//Protected Endpoints

	r.Handle("/people", authentication.ValidateToken.Handler(handler.GetPeopleEndPointHandler)).Methods("GET")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(handler.GetPersonEndPointHandler)).Methods("GET")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(handler.CreatePersonEndPointHandler)).Methods("POST")
	r.Handle("/people/{id}", authentication.ValidateToken.Handler(handler.DeletePersonEndPointHandler)).Methods("DELETE")

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

func main() {
	authentication.InitKeys()
	data.InitData(data.GetPeople())
	StartServer()
}
