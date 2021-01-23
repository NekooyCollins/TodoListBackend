package main

import (
	"./database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Cache active group task on server.
type groupTask struct {
	taskInfo   database.TaskType
	member     map[string]bool
}

var cacheGroupTaskList []groupTask

// Handle group task check.
func getGroupTaskState(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	// Get user id
	inputID := r.URL.Query().Get("id")
	fmt.Println("Ask for user data of:", inputID)

	// Check from cache group task list.
	for _, task := range cacheGroupTaskList {
		if val, ok := task.member[inputID]; ok{
			if val == false {
				task.member[inputID] = true
				taskJson, err := json.Marshal(task.taskInfo)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(taskJson)
				return
			}
		}
	}

	// If not in group task, return StatusBadRequest.
	http.Error(w, "Not being invited to group task", http.StatusBadRequest)
	return
}

// A member in a group task start the task.
func postStartGroupTask(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()

	// Parse json
	var formData database.TaskType
	json.NewDecoder(r.Body).Decode(&formData)

	var newGroupTask groupTask
	newGroupTask.taskInfo = formData

	// Find members of this task.
	// Get user tasks.
	queryStr := "SELECT * FROM usertasklist WHERE taskid=" + strconv.Itoa(newGroupTask.taskInfo.ID) + ";"
	gettask, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for gettask.Next() {
		var tmp database.UserTaskListType
		if err = gettask.Scan(&tmp.UserID, &tmp.TaskID); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		newGroupTask.member[strconv.Itoa(tmp.UserID)] = false
	}

	// If already in a group task, retrun with StatusBadGateway.
	for _, item := range cacheGroupTaskList {
		if item.taskInfo.ID == newGroupTask.taskInfo.ID{
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	cacheGroupTaskList = append(cacheGroupTaskList, newGroupTask)
	w.WriteHeader(http.StatusOK)
	return
}

// Accept invitation and join group task.
func joinGroupTask(w http.ResponseWriter, r *http.Request){
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputUserID,_ := formData["userid"]
	inputTaskID,_ := strconv.Atoi(formData["taskid"])

	for _, item := range cacheGroupTaskList{
		if item.taskInfo.ID == inputTaskID{
			if _, ok := item.member[inputUserID]; ok{
				item.member[inputUserID] = true
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	http.Error(w, "Failed to join the task.", http.StatusBadRequest)
	return
}

// Check if all members have joined one group task.
func checkStartGroupTask(w http.ResponseWriter, r *http.Request){
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID,_ := strconv.Atoi(formData["taskid"])

	for _, item := range cacheGroupTaskList{
		if item.taskInfo.ID == inputTaskID{
			for _, val := range item.member{
				if val == false {
					http.Error(w, "Failed to start the task.", http.StatusBadRequest)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Failed to start the task.", http.StatusBadRequest)
	return
}

// Quit a group task.
func quitGroupTask(w http.ResponseWriter, r *http.Request){
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID,_ := strconv.Atoi(formData["taskid"])

	for _, item := range cacheGroupTaskList{
		if item.taskInfo.ID == inputTaskID{
			for key, _ := range item.member{
				item.member[key] = false
			}
		}
	}

	http.Error(w, "Failed to join the task.", http.StatusBadRequest)
	return
}

// Periodically check if an on-going group task is quit by others,
// if quit, response 200.
func checkGroupTaskQuit(w http.ResponseWriter, r *http.Request){
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID,_ := strconv.Atoi(formData["taskid"])

	for _, item := range cacheGroupTaskList{
		if item.taskInfo.ID == inputTaskID{
			for _, val := range item.member{
				if val == false {
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}
	}

	http.Error(w, "Failed to start the task.", http.StatusBadRequest)
	return
}
