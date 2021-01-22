package main

import (
	"./database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type addTaskStruct struct {
	ID          int                 `json:"id"`
	CreatorID   int                 `json:"creatorid"`
	Title       string              `json:"title"`
	Desc        string              `json:"description"`
	Duration    int                 `json:"duration"`
	RemainTime  int                 `json:"remaintime"`
	Type        string              `json:"typestr"`
	IsFinish    bool                `json:"isfinish"`
	IsGroupTask bool                `json:"isgrouptask"`
	Members     []string		 	`json:"member"`
}

// Get email and password from login view
// and verify.
func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()

	// Parse json
	var formData addTaskStruct
	json.NewDecoder(r.Body).Decode(&formData)

	// Check if input data is legal
	if formData.Title == "" {
		http.Error(w, "All fields can't be empty.", http.StatusBadRequest)
		return
	}

	fmt.Println("Get add task title: ", formData.Title)
	fmt.Println("Description:", formData.Desc)

	var isGrouptask string = "false"

	// Add task into 'task' table
	insertSql := "INSERT INTO task(title, descption, duration, remaintime, typestr, isfinish, isgrouptask) VALUES ('" +
		formData.Title + "', '" + formData.Desc + "', " + strconv.Itoa(formData.Duration) + ", " + strconv.Itoa(formData.Duration) +
		", '" + formData.Type + "', false, " + isGrouptask + ")"

	res, err := dbconn.DBConn.Exec(insertSql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskID, err := res.LastInsertId()

	// Add user-task relationship into usertasklist table
	insertSql = "INSERT INTO usertasklist VALUES (" + strconv.Itoa(formData.CreatorID) + ", " + strconv.FormatInt(taskID, 10) + ")"
	_, err = dbconn.DBConn.Exec(insertSql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, userEmail := range formData.Members {
		queryStr := "SELECT * FROM user WHERE email='" + userEmail + "';"
		rets, err := dbconn.DBConn.Query(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userdata database.UserType
		// Get user data
		for rets.Next() {
			if err = rets.Scan(&userdata.ID, &userdata.Name, &userdata.Email, &userdata.Passwd); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
		}

		insertSql = "INSERT INTO usertasklist VALUES (" + strconv.Itoa(userdata.ID) + ", " + strconv.FormatInt(taskID, 10) + ")"
		_, err = dbconn.DBConn.Exec(insertSql)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update task to group task
		updateSql := "UPDATE task SET isgrouptask=true WHERE id=" + strconv.FormatInt(taskID, 10);
		_, err = dbconn.DBConn.Exec(updateSql)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Successful created user. Task id is"+strconv.FormatInt(taskID, 10))
}
