package main

import database "./database"

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
	switch r.Method {
	case "GET":
		rets, err := dbconn.DBConn.Query("SELECT * FROM user")
		if err!=nil{
			panic(err.Error())
		}

		for rets.Next(){
			var userdata database.UserType
			if err = rets.Scan(&userdata.ID, &userdata.Email, &userdata.Name, &userdata.Passwd);err != nil{
				panic(err.Error())
			}

			fmt.Printf(string(userdata.ID), userdata.Name, userdata.Email)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
