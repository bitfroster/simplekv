package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db = NewDb("simple")

type Payload struct {
	Dbkey   string
	Dbvalue string
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/db", CreateRecord).Methods("POST")
	router.HandleFunc("/db/{recordId}", GetRecord).Methods("GET")
	router.HandleFunc("/db/{recordId}", DeleteRecord).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["recordId"]
	result, exists := db.Get(id)

	if !exists {
		respondWithError(w, http.StatusBadRequest, "Invalid record ID")
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

func CreateRecord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Payload
	err := decoder.Decode(&p)

	if err != nil {
		panic(err)
	}

	db.Set(p.Dbkey, []byte(p.Dbvalue))

	respondWithJSON(w, http.StatusOK, p)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["recordId"]
	result, exists := db.Delete(id)

	if !exists {
		respondWithError(w, http.StatusBadRequest, "Invalid record ID")
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
