package main

import (
	"./database"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

// Handle getranklist request,
// return all focustime of friends.
func getRankList(w http.ResponseWriter, r *http.Request) {
	// Return list
	var retRankList []database.RankListType

	// Get request query value.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	inputUserID := r.URL.Query().Get("userid")
	if inputUserID == "" {
		http.Error(w, "Can't get value.", http.StatusBadRequest)
		return
	}

	// Get user's total focus time
	queryStr := "SELECT * FROM user WHERE id='" + inputUserID+ "';"
	userRet, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for userRet.Next(){
		var userInfo database.UserType
		var myFinishedTaskList []database.TaskType
		var myRankItem database.RankListType
		var myFinishedTaskSetEmpty bool = true
		var myTotalFocusTime int = 0

		if err = userRet.Scan(&userInfo.ID, &userInfo.Name, &userInfo.Email, &userInfo.Passwd); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		// Get user's task id
		queryStr = "SELECT * FROM usertasklist WHERE userid='" + userInfo.ID + "';"
		myTaskIdRet, err := dbconn.DBConn.Query(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for myTaskIdRet.Next() {
			var userAndTask database.UserTaskListType
			var finishedTask database.TaskType
			if err = myTaskIdRet.Scan(&userAndTask.UserID, &userAndTask.TaskID); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}

			// Get finished task set from task table
			queryStr = "SELECT * FROM task WHERE id='" + userAndTask.TaskID + "' AND isfinish=1" + ";"
			myTasksRet, err := dbconn.DBConn.Query(queryStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			for myTasksRet.Next() {
				myFinishedTaskSetEmpty = false
				if err = myTasksRet.Scan(&finishedTask.ID, &finishedTask.Title, &finishedTask.Desc, &finishedTask.Duration,
					&finishedTask.RemainTime, &finishedTask.Type, &finishedTask.IsFinish, &finishedTask.IsGroupTask); err != nil {
					http.Error(w, err.Error(), http.StatusBadGateway)
				}
				myFinishedTaskList = append(myFinishedTaskList, finishedTask)
			}
		}
		if myFinishedTaskSetEmpty == true {
			myTotalFocusTime = 0
		} else {
			for i := 0; i < len(myFinishedTaskList); i++ {
				myTotalFocusTime = myTotalFocusTime + myFinishedTaskList[i].Duration
			}
		}
		myRankItem.UserName = userInfo.Name
		myRankItem.TotalFocusTime = myTotalFocusTime
		retRankList = append(retRankList, myRankItem)
	}

	// Get all friend's ID from userfriendlist.
	queryStr = "SELECT * FROM userfriendlist WHERE userid='" + inputUserID + "';"
	friendRets, err := dbconn.DBConn.Query(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for friendRets.Next() {
		// Variables for one friend
		var userAndFriend database.UserFriendListType
		var friendInfo database.UserType
		var totalTimeOfOneFriend int = 0
		var singleRankItem database.RankListType
		var finishedTaskList []database.TaskType
		var isTaskEmpty bool = true

		if err = friendRets.Scan(&userAndFriend.UserID, &userAndFriend.FriendID); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		// Get this friends' all info from user.
		queryStr = "SELECT * FROM user WHERE id='" + userAndFriend.FriendID + "';"
		friendInfoRet, err := dbconn.DBConn.Query(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Get friend's ID and Name
		for friendInfoRet.Next() {
			if err = friendInfoRet.Scan(&friendInfo.ID, &friendInfo.Name, &friendInfo.Email, &friendInfo.Passwd); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}

			// Get all taskID from usertasklist table
			queryStr = "SELECT * FROM usertasklist WHERE userid='" + friendInfo.ID + "';"
			taskIdRet, err := dbconn.DBConn.Query(queryStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			for taskIdRet.Next() {
				isTaskEmpty = false
				var userAndTask database.UserTaskListType
				var finishedTask database.TaskType
				if err = taskIdRet.Scan(&userAndTask.UserID, &userAndTask.TaskID); err != nil {
					http.Error(w, err.Error(), http.StatusBadGateway)
				}

				// Get info of this task from task table
				queryStr = "SELECT * FROM task WHERE id='" + userAndTask.TaskID + "' AND isfinish=1" + ";"
				tasksRet, err := dbconn.DBConn.Query(queryStr)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				for tasksRet.Next() {
					if err = tasksRet.Scan(&finishedTask.ID, &finishedTask.Title, &finishedTask.Desc, &finishedTask.Duration,
						&finishedTask.RemainTime, &finishedTask.Type, &finishedTask.IsFinish, &finishedTask.IsGroupTask); err != nil {
						http.Error(w, err.Error(), http.StatusBadGateway)
					}
				}
				finishedTaskList = append(finishedTaskList, finishedTask)
			}
			// Get total-focus-time of one friend
			singleRankItem.UserName = friendInfo.Name
			if isTaskEmpty == true {
				totalTimeOfOneFriend = 0
			} else {
				for i := 0; i < len(finishedTaskList); i++ {
					totalTimeOfOneFriend = totalTimeOfOneFriend + finishedTaskList[i].Duration
				}
			}
			singleRankItem.TotalFocusTime = totalTimeOfOneFriend
			retRankList = append(retRankList, singleRankItem)
		}
	}

	sort.Slice(retRankList, func(i, j int) bool {
		if retRankList[i].TotalFocusTime > retRankList[j].TotalFocusTime {
			return true
		}
		return false
	})

	for i := 0; i < len(retRankList); i++ {
		fmt.Println("now is " + retRankList[i].UserName)
		fmt.Println("his total focus time is " + strconv.Itoa(retRankList[i].TotalFocusTime))
	}

	// Return json data.
	friendListJson, err := json.Marshal(retRankList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(friendListJson)
}
