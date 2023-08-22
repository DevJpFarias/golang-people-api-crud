package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	routes := mux.NewRouter().StrictSlash(true)

	routes.HandleFunc("/persons", getAll).Methods("GET")
	routes.HandleFunc("/persons", create).Methods("POST")
	var port = ":3000"

	fmt.Println("Server Running on port:", port)
	log.Fatal(http.ListenAndServe(port, routes))
}

type Person struct {
	Name string
}

var people = []Person{
	{
		Name: "Fulano",
	},
	{
		Name: "Fulaninho",
	},
}

func getAll(response http.ResponseWriter, request *http.Request) {
	var validPeople = filterValidNames(people)

	json.NewEncoder(response).Encode(validPeople)
}

func create(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var person Person

	body, err := io.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}

	if err := request.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &person); err != nil {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")

		response.WriteHeader(422)

		if err := json.NewEncoder(response).Encode(err); err != nil {
			panic(err)
		}
	}

	json.Unmarshal(body, &person)

	people = append(people, person)

	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(response).Encode(person); err != nil {
		panic(err)
	}
}

func filterValidNames(people []Person) []Person {
	var validPeople []Person

	for _, person := range people {
		if person.Name != "" {
			validPeople = append(validPeople, person)
		}
	}

	return validPeople
}
