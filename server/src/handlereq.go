package main

import (
	"./database"
	"encoding/json"
	"strconv"
)

import (
	"fmt"
	"log"
	"net/http"
)

func handleRequests() {
	http.HandleFunc("/login", loginCheck)
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/getuserdata", getUserData)
	http.HandleFunc("/gettasklist", getTaskList)
	http.HandleFunc("/gettaskmember", getTaskMember)
	http.HandleFunc("/gettaskdetail", getTaskDetail)
	http.HandleFunc("/addtask", addTask)
	http.HandleFunc("/getfriendlist", getFriendList)
	http.HandleFunc("/getranklist", getRankList)
	http.HandleFunc("/addfriend", addFriend)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Get email and password from login view
// and verify.
func loginCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	// Parse json
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)

	inputEmail := formData["email"]
	inputPasswd := formData["passwd"]

	// check from database
	rets, err := dbconn.DBConn.Query("SELECT * FROM user WHERE email='" + inputEmail + "';")
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if rets == nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Verify user email and password
	for rets.Next() {
		var userdata database.UserType
		if err = rets.Scan(&userdata.ID, &userdata.Name, &userdata.Email, &userdata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		if userdata.Email == inputEmail && userdata.Passwd == inputPasswd {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Wrong email or password", http.StatusForbidden)
}

// Get new user data to register.
func registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	// Parse input data
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputUser := formData["name"]
	inputEmail := formData["email"]
	inputPasswd := formData["passwd"]

	// Check if input data is legal
	if inputUser == "" || inputPasswd == "" || inputEmail == "" {
		http.Error(w, "All fields can't be empty.", http.StatusBadRequest)
		return
	}

	// check from database if email has existed
	var retCnt int
	_ = dbconn.DBConn.QueryRow("SELECT count(id) FROM user WHERE email='" + inputEmail + "';").Scan(&retCnt)
	if retCnt != 0 {
		http.Error(w, "User has existed.", http.StatusBadRequest)
		return
	}

	// insert into database
	insertSql := "INSERT INTO user(name, email, passwd) VALUES ('" + inputUser + "', '" + inputEmail + "', '" + inputPasswd + "')"

	res, err := dbconn.DBConn.Exec(insertSql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lastID, err := res.LastInsertId()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Successful created user. User id is"+string(lastID))
	return
}

// Handle getuserdata get request,
// return all user data of one user.
func getUserData(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputEmail := r.URL.Query().Get("email")
	fmt.Println("Ask for user data of:", inputEmail)
	if inputEmail == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Check from database.
	queryStr := "SELECT * FROM user WHERE email='" + inputEmail + "';"
	rets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send back json data.
	for rets.Next() {
		var userdata database.UserType
		if err = rets.Scan(&userdata.ID, &userdata.Name, &userdata.Email, &userdata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		userJson, err := json.Marshal(userdata)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	}
}

// Handle gettasklist request,
// return all tasks list of the user.
func getTaskList(w http.ResponseWriter, r *http.Request) {
	var userdata database.UserType
	var retTaskList []database.TaskType

	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputEmail := r.URL.Query().Get("email")
	if inputEmail == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Check from database.
	queryStr := "SELECT * FROM user WHERE email='" + inputEmail + "';"
	rets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get user id.
	for rets.Next() {
		if err = rets.Scan(&userdata.ID, &userdata.Email, &userdata.Name, &userdata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		break
	}
	// Get user tasks.
	queryStr = "SELECT * FROM usertasklist WHERE userid=" + strconv.Itoa(userdata.ID) + ";"
	gettask, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for gettask.Next() {
		var usertask database.UserTaskListType
		var task database.TaskType
		if err = gettask.Scan(&usertask.UserID, &usertask.TaskID); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		queryStr = "SELECT * FROM task WHERE id=" + strconv.Itoa(usertask.TaskID) + ";"
		retTask, err := dbconn.DBConn.Query(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for retTask.Next() {
			if err = retTask.Scan(&task.ID, &task.Title, &task.Desc, &task.Duration, &task.RemainTime, &task.Type,
				&task.IsFinish, &task.IsGroupTask); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			retTaskList = append(retTaskList, task)
		}
	}

	// Return json data.
	taskListJson, err := json.Marshal(retTaskList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(taskListJson)
}

// Return member array of one task
func getTaskMember(w http.ResponseWriter, r *http.Request) {
	var retUserList []database.UserType

	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputTaskID := r.URL.Query().Get("taskid")
	if inputTaskID == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Check from database.
	queryStr := "SELECT * FROM usertasklist WHERE taskid='" + inputTaskID + "';"
	rets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get user id.
	for rets.Next() {
		var getTask database.UserTaskListType
		if err = rets.Scan(&getTask.UserID, &getTask.TaskID); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		queryUserStr := "SELECT * FROM user WHERE id=" + strconv.Itoa(getTask.UserID) + ";"
		userRet, err := dbconn.DBConn.Query(queryUserStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for userRet.Next() {
			var userItem database.UserType
			if err = userRet.Scan(&userItem.ID, &userItem.Name, &userItem.Email, &userItem.Passwd); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			retUserList = append(retUserList, userItem)
		}
	}

	// Return json data.
	userListJson, err := json.Marshal(retUserList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userListJson)
}

// Return task detail to client.
func getTaskDetail(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputTaskID := r.URL.Query().Get("taskid")
	if inputTaskID == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Check from database.
	queryStr := "SELECT * FROM task WHERE id=" + inputTaskID + ";"
	rets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send back json data.
	for rets.Next() {
		var taskData database.TaskType
		if err = rets.Scan(&taskData.ID, &taskData.Title, &taskData.Desc, &taskData.Duration, &taskData.RemainTime,
			&taskData.Type, &taskData.IsFinish, &taskData.IsGroupTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		userJson, err := json.Marshal(taskData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	}
}

// Handle getfriendlist request,
// return all friends of the user.
func getFriendList(w http.ResponseWriter, r *http.Request) {
	var userdata database.UserType
	var retFriendList []database.UserType

	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputEmail:= r.URL.Query().Get("email")
	if inputEmail == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Check from database.
	queryStr := "SELECT * FROM user WHERE email='"+inputEmail+"';"
	rets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get user id.
	for rets.Next() {
		if err = rets.Scan(&userdata.ID, &userdata.Email, &userdata.Name, &userdata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		break
	}
	println("now user name is: " + userdata.Name )
	println("now userID is: " + strconv.Itoa(userdata.ID))
	// Get user friends.
	queryStr = "SELECT * FROM userfriendlist WHERE userid="+strconv.Itoa(userdata.ID)+";"
	getfriend, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for getfriend.Next() {
		var userfriends database.UserFriendListType
		var friend database.UserType
		if err = getfriend.Scan(&userfriends.UserID, &userfriends.FriendID); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		queryStr = "SELECT * FROM user WHERE id="+strconv.Itoa(userfriends.FriendID)+";"
		retFriends, err := dbconn.DBConn.Query(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for retFriends.Next(){
			if err = retFriends.Scan(&friend.ID, &friend.Name, &friend.Email, &friend.Passwd); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			retFriendList = append(retFriendList, friend)
		}
		println("now friend name is: " + friend.Name )
		println("now friendID is: " + strconv.Itoa(friend.ID))
	}

	// Return json data.
	friendListJson, err := json.Marshal(retFriendList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(friendListJson)
}

// insert friend relationship
func addFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()

	// Parse json
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)

	myEmail := formData["myemail"]
	friendEmail := formData["friendemail"]

	// get userID
	userRet, err:= dbconn.DBConn.Query("SELECT * FROM user WHERE email='"+myEmail+"';")
	if err != nil{
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if userRet == nil{
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// get friend's ID
	friendRet, err:= dbconn.DBConn.Query("SELECT * FROM user WHERE email='"+friendEmail+"';")
	if err != nil{
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if friendRet == nil{
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var userdata database.UserType
	var frienddata database.UserType

	for userRet.Next(){
		if err = userRet.Scan(&userdata.ID, &userdata.Name, &userdata.Email, &userdata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
	for friendRet.Next(){
		if err = friendRet.Scan(&frienddata.ID, &frienddata.Name, &frienddata.Email, &frienddata.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
	// Add user-friend relationship into usertasklist table
	insertSql := "INSERT INTO userfriendlist VALUES ("+strconv.Itoa(userdata.ID)+", "+strconv.Itoa(frienddata.ID)+")"
	_, err = dbconn.DBConn.Exec(insertSql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
