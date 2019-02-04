package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

var session, err = mgo.Dial("http://mongodbservice:27017")

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!\n"))
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	err = session.DB("PruebaDB").C("person").Find(nil).All(&result)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(result)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err = session.DB("PruebaDB").C("person").Find(bson.M{"id": id}).All(&resultId)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(resultId)
}

func GetPersonByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nm := vars["nm"]

	err = session.DB("PruebaDB").C("person").Find(bson.M{"firstname": nm}).All(&resultId)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(resultId)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
        firstnm := vars["firstnm"]
	lastnm := vars["lastnm"]

	c := session.DB("PruebaDB").C("person")
	err = c.Insert(&Person{Id: id, Firstname: firstnm, Lastname: lastnm, Address: &Address{City: "City X", State: "State X"}})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("Persona Creada!\n"))
}

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

func main() {
	//c := session.DB("PruebaDB").C("person")
	//err = c.Insert(&Person{Id: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}},
	//	&Person{Id: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})

	router := mux.NewRouter()
	router.HandleFunc("/", Welcome).Methods("GET")
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/byname/{nm}", GetPersonByName).Methods("GET")
	router.HandleFunc("/people/add/{id}/{firstnm}/{lastnm}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/upd/{id}/{firstnm}/{lastnm}", UpdatePerson).Methods("PUT")
	router.HandleFunc("/people/del/{id}", DeletePerson).Methods("DELETE")
	http.ListenAndServe(":8001", router)
}