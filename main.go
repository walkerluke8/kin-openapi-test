// 2020
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type Server struct {
	db *cache.Cache
}

func (s Server) GetPets(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	items := s.db.Items()

	pets := make([]Pet, 0, len(items))
	for _, tx := range items {
		pets = append(pets, tx.Object.(Pet))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets)
}

func (s Server) GetPet(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	petID := vars["petID"]
	result := Pet{}

	foo, found := s.db.Get(petID)
	if found {
		result = foo.(Pet)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (s Server) AddPet(w http.ResponseWriter, r *http.Request) {

	var p Pet

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	s.db.Set(strconv.Itoa(p.ID), p, cache.DefaultExpiration)
	if err != nil {
		panic(err)
	}
}

func main() {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	server := Server{
		db: c,
	}

	r := mux.NewRouter()
	r.HandleFunc("/pets", server.GetPets).Methods("GET")
	r.HandleFunc("/pets", server.AddPet).Methods("POST")
	r.HandleFunc("/pets/{petID}", server.GetPet)

	logrus.Fatal(http.ListenAndServe(":8000", r))
}
