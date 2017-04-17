package data

import "fmt"

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

func GetPeople() []Person {
	return people
}

func SetPeople(people []Person) {
	fmt.Println("Yet to implement this method: Set People")
}
func InitData(people []Person) {
	people = append(people, Person{ID: "1", Firstname: "Ashwini", Lastname: "Patankar", Address: &Address{City: "Bangalore", State: "India"}})
	people = append(people, Person{ID: "2", Firstname: "Manish", Lastname: "Patankar", Address: &Address{City: "San Fransico", State: "California"}})
	people = append(people, Person{ID: "3", Firstname: "Hun", Lastname: "Patankar", Address: &Address{City: "Munich", State: "Germany"}})
}
