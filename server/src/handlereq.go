package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/getuserdata", getUserData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")

}
