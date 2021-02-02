package main

import (
	"./database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Cache active group task on server.
type groupTask struct {
	taskInfo  database.TaskType
	member    map[string]bool
	alertFlag map[string]bool
	quitFlag  bool
}

var cacheGroupTaskList []groupTask

// Handle group task check.
func getGroupTaskState(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	//fmt.Println("get group task state is called!")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	// Get user id
	inputID := r.URL.Query().Get("id")
	inputID = strings.ToLower(inputID)
	fmt.Println("input ID is: ", inputID)

	// Check from cache group task list.
	for _, task := range cacheGroupTaskList {
		//fmt.Println("cached group task is:")
		//fmt.Println(cacheGroupTaskList)
		if val, ok := task.alertFlag[inputID]; ok {
			if val == false {
				// If not join yet, alert
				task.alertFlag[inputID] = true
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
	//fmt.Println("Try to start a group task.")
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
	newGroupTask.member = make(map[string]bool)
	newGroupTask.alertFlag = make(map[string]bool)

	// Find members of this task.
	// Get user tasks.
	queryStr := "SELECT * FROM usertasklist WHERE taskid='" + newGroupTask.taskInfo.ID + "';"
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
		tmp.TaskID = strings.ToLower(tmp.TaskID)
		tmp.UserID = strings.ToLower(tmp.UserID)
		newGroupTask.member[tmp.UserID] = false
		newGroupTask.alertFlag[tmp.UserID] = false
		newGroupTask.quitFlag = false
	}

	// If already in a group task, retrun with StatusBadGateway.
	for _, item := range cacheGroupTaskList {
		if item.taskInfo.ID == newGroupTask.taskInfo.ID {
			fmt.Println("Some other member already started a task")
			http.Error(w, "Some other member already started a task", http.StatusBadGateway)
			return
		}
	}

	cacheGroupTaskList = append(cacheGroupTaskList, newGroupTask)
	w.WriteHeader(http.StatusOK)
	fmt.Println("Group task ", newGroupTask.taskInfo.Title, "has started.")
	return
}

// Accept invitation and join group task.
func joinGroupTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("join the task is called!!!!!!!!!!!!!")
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputUserID, _ := formData["userid"]
	inputTaskID, _ := formData["taskid"]
	inputUserID = strings.ToLower(inputUserID)
	inputTaskID = strings.ToLower(inputTaskID)
	//fmt.Println("User", inputUserID, "wants to join task", inputTaskID)

	for idx, item := range cacheGroupTaskList {
		if strings.ToLower(item.taskInfo.ID) == inputTaskID {
			for k, _ := range item.member {
				//fmt.Println("here K is:", k)
				if k == inputUserID{
					cacheGroupTaskList[idx].member[inputUserID] = true
					cacheGroupTaskList[idx].alertFlag[inputUserID] = true
					w.WriteHeader(http.StatusOK)
					fmt.Println(inputUserID, "successfully join the task", inputTaskID)
					return
				}
			}
		}
	}

	fmt.Println("join group task: ", inputUserID, "failed join the task", inputTaskID)
	http.Error(w, "Failed to join the task.", http.StatusBadRequest)
	return
}

// Check if all members have joined one group task.
func checkStartGroupTask(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	//fmt.Println("check start group task is called!!!!")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID := formData["taskid"]
	inputTaskID = strings.ToLower(inputTaskID)
	//fmt.Println("Check if task", inputTaskID, "could start.")

	for _, item := range cacheGroupTaskList {
		if strings.ToLower(item.taskInfo.ID) == inputTaskID {
			for _, val := range item.member {
				if val == false {
					fmt.Println("check start group task: Failed to start the task.")
					//fmt.Println(cacheGroupTaskList)
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
func quitGroupTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("quit group task is called!!!!")
	// Get request query value.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID := formData["taskid"]
	inputTaskID = strings.ToLower(inputTaskID)
	fmt.Println("task ", inputTaskID, "want to quit.")

	for idx, item := range cacheGroupTaskList {
		if strings.ToLower(item.taskInfo.ID) == inputTaskID {
			cacheGroupTaskList[idx].quitFlag = true
			fmt.Println("Task ", inputTaskID, "has quit.")
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Failed to quit the task.", http.StatusBadRequest)
	return
}

// Periodically check if an on-going group task is quit by others,
// if quit, response 200.
func checkGroupTaskQuit(w http.ResponseWriter, r *http.Request) {
	// Get request query value.
	fmt.Println("Check for task quit status here!!!!!!!!")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	r.ParseForm()
	formData := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&formData)
	inputTaskID := formData["taskid"]
	inputTaskID = strings.ToLower(inputTaskID)
	fmt.Println("check if task ", inputTaskID, "quit.")

	for _, item := range cacheGroupTaskList {
		if strings.ToLower(item.taskInfo.ID) == inputTaskID {
			if item.quitFlag == true {
				fmt.Println("Found task ", inputTaskID, "has quit.")
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	http.Error(w, "Task isn't quited.", http.StatusBadRequest)
	return
}
