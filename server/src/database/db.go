package database

import (
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// TaskType receives original task data from db
type TaskType struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Desc        string `json:"descption"`
	Duration    int    `json:"duration"`
	RemainTime  int    `json:"remaintime"`
	Type        string `json:"typestr"`
	IsFinish    bool   `json:"isfinish"`
	IsGroupTask bool   `json:"isgrouptask"`
}

// UserType receives original user data from db
type UserType struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

// UserTaskListType receives original user-task
// list data from db
type UserTaskListType struct {
	UserID int `json:"userid"`
	TaskID int `json:"taskid"`
}

// UserFriendList receives original user-friend
// list data from DB
type UserFriendListType struct {
	UserID   int `json:"userid"`
	FriendID int `json:"friendid"`
}

// DBType for chain call
type DBType struct {
	DBConn *sql.DB
	Error  error
}

// DBConnect function opens and returns a connection to mysql DB
func (db *DBType) DBConnect() *DBType {
	DBConfig, err := GetDBConfig()
	if err != nil {
		panic(err.Error())
	}
	connStr := DBConfig.DBUsername + ":" +
		DBConfig.DBPasswd + "@tcp(" +
		DBConfig.DBHost + ":" +
		DBConfig.DBPort + ")/" +
		DBConfig.DBName
	fmt.Printf("DB address "+connStr)
	db.DBConn, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connection to database successful.")

	db.DBConn.SetMaxIdleConns(64)
	db.DBConn.SetMaxOpenConns(64)
	db.DBConn.SetConnMaxLifetime(time.Minute)
	return db
}

// DBClose close databse connection
func (db *DBType) DBClose() {
	db.DBConn.Close()
}
