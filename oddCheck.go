package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Value struct {
	value int  `json:"value"`
	isodd bool `json:"isodd"`
}

var values []Value

func main() {
	values = append(values, Value{123, false})
	values = append(values, Value{122, true})
	values = append(values, Value{111, false})
	mux := mux.NewRouter()
	mux.HandleFunc("/", check)
	mux.HandleFunc("/getvalues", getValues)
	mux.HandleFunc("/getvalue/{value}", getValue)
	mux.HandleFunc("/createvalue/{value}", createValue)
	mux.HandleFunc("/deletevalue/{value}", deleteValue)

	log.Print("Идёт подклчюение к http://localhost:4040")
	err := http.ListenAndServe("localhost:4040", mux)
	if err != nil {
		return
	}
}
func check(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("Connected"))
}
func getValues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte("Used wrong call method. Try to use GET"))
		return
	}
	if len(values) == 0 {
		w.Write([]byte("List is empty \n"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(values)
	}
}
func getValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Write([]byte("Used wrong call method. Try to use GET"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range values {
		if string(item.value) == params[`value`] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Value{})
	w.Write([]byte("Value not founded"))
}

func createValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("Wrong use. Try to use method POST"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var value Value
	_ = json.NewDecoder(r.Body).Decode(&value)
	if value.value%2 == 0 {
		value.isodd = true
	} else {
		value.isodd = false
	}
	values = append(values, value)
	json.NewEncoder(w).Encode(value)

}

func deleteValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Fatal("Wrong use. Try to use method DELETE")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range values {
		if string(item.value) == params[`value`] {
			values = append(values[:index], values[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(values)
}
