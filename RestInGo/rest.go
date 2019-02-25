package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type Persons struct {
	Persons []Person `json:"Persons"`
}

type Person struct {
	ID string `json:"ID"`
	Name string `json:"Name"`
	Surname string `json:"Surname"`
	Age int `json:"Age"`

}

var allPersons = listOfUsers().Persons

func handler(w http.ResponseWriter, r *http.Request) {
	jsonVar, err := os.Open("json/package.json")
	if err != nil{
		fmt.Println(err)
	}
	defer jsonVar.Close()
	byteValue, _  := ioutil.ReadAll(jsonVar)

	var persons Persons
	_ = json.Unmarshal(byteValue, &persons)

	tmpl, err := template.ParseFiles("html/hello.html")

	_ = tmpl.ExecuteTemplate(w, "jsonfile", persons)
}

func listOfUsers() Persons{
	jsonVar, err := os.Open("json/package.json")
	if err != nil{
		fmt.Println(err)
	}
	defer jsonVar.Close()
	byteValue, _  := ioutil.ReadAll(jsonVar)

	var person Persons
	_ = json.Unmarshal(byteValue, &person)

	return person
}




func getPersons(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(listOfUsers())
}


func getPerson(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range listOfUsers().Persons{
		if item.ID == params["id"]{
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createPerson(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	var newPerson Person

	_ = json.NewDecoder(r.Body).Decode(&newPerson)
	newPerson.ID = strconv.Itoa(len(allPersons)+1)
	allPersons = append(allPersons,newPerson)
	_ = json.NewEncoder(w).Encode(allPersons)
}

func updatePersons(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newPerson Person
	_ = json.NewDecoder(r.Body).Decode(&newPerson)
	for index, item := range allPersons{
		if item.ID == params["id"]{
			newPerson.ID = params["id"]
			allPersons[index] = newPerson
		}
	}
	_ = json.NewEncoder(w).Encode(allPersons)
}

func deletePersons(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range allPersons{
		if item.ID == params["id"]{
			allPersons = append(allPersons[:index], allPersons[index + 1:]...)
			break
		}
	}
	_ = json.NewEncoder(w).Encode(allPersons)
}



func main() {
	router := mux.NewRouter()

	router.HandleFunc("/person", getPersons).Methods("Get")
	router.HandleFunc("/person/{id}", getPerson).Methods("Get")
	router.HandleFunc("/person", createPerson).Methods("POST")
	router.HandleFunc("/person/{id}", updatePersons).Methods("PUT")
	router.HandleFunc("/person/{id}", deletePersons).Methods("DELETE")

	http.HandleFunc("/firstTask", handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
