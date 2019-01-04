package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The person Type (more like an object)
type Person struct {
	Id        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person
var result []Person
var resultId []Person

var session, err = mgo.Dial("localhost")

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!\n"))
}

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	err = session.DB("PruebaDB").C("person").Find(nil).All(&result)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err = session.DB("PruebaDB").C("person").Find(bson.M{"id": id}).All(&resultId)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(resultId)
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	c := session.DB("PruebaDB").C("person")
	err = c.Insert(&Person{Id: id, Firstname: "Create", Lastname: "Create", Address: &Address{City: "City X", State: "State X"}})

	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte("Persona Creada!\n"))
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err = session.DB("PruebaDB").C("person").Remove(bson.M{"id": id})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("Persona Borrada!\n"))
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	firstnm := vars["firstnm"]
	lastnm := vars["lastnm"]

	err = session.DB("PruebaDB").C("person").Update(bson.M{"id": id}, bson.M{"$set": bson.M{"firstname": firstnm, "lastname": lastnm}})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("Persona Actualizada!\n"))
}

// main function to boot up everything
func main() {
	c := session.DB("PruebaDB").C("person")
	err = c.Insert(&Person{Id: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}},
		&Person{Id: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", Welcome).Methods("GET")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/add/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/upd/{id}/{firstnm}/{lastnm}", UpdatePerson).Methods("PUT")
	router.HandleFunc("/people/del/{id}", DeletePerson).Methods("DELETE")
	http.ListenAndServe(":8001", router)
}
