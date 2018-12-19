package main

import (
	"fmt"
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
	"gopkg.in/mgo.v2"

)

// The person Type (more like an object)
type Person struct {
    ID        string
    Firstname string
    Lastname  string
    Address   *Address
}
type Address struct {
    City  string
    State string
}

var people []Person
var result []Person

func Welcome(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome!\n"))
}

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
    //w.Write([]byte(result))
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var person Person
    _ = json.NewDecoder(r.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
        json.NewEncoder(w).Encode(people)
    }
}

// main function to boot up everything
func main() {

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("PruebaDB").C("person")
	err = c.Insert(&Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}},
	&Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	if err != nil {
			log.Fatal(err)
	}


		
	err = session.DB("PruebaDB").C("person").Find(nil).All(&result)
    if err != nil {
    	log.Fatal(err)
	}

	fmt.Println("Person: ", result)

	
    router := mux.NewRouter()
    router.HandleFunc("/", Welcome).Methods("GET")
    router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8001", router))
}