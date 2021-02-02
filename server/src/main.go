package main

import (
	database "./database"
)

var dbconn database.DBType

func main() {
	// Open database connection
	dbconn.DBConnect()
	defer dbconn.DBClose()

	handleRequests()
}
