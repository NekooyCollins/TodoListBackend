# TodoListBackend
The backend server for TodoList app implemented by Golang

## functions
* handle request from client
* update database in server

## DataBase Tables
### Task

| Attribute | Type |
| --------- | ---- |
| uuid  | uuid |
| title | string |
| descption | string |
| duration | int |
| remaintime | int |
| typestr | string |
| isfinish | bool |
| isgrouptask | bool |

### User

| Attribute | Type |
| --------- | ---- |
| uuid | uuid |
| name | string |
| email | string |
| passwd | string |

### UserTaskList

| Attribute | Type |
| --------- | ---- |
| userid | uuid |
| taskid | uuid |

### UserFriendList

| Attribute | Type |
| --------- | ---- |
| userid | uuid |
| friendid| uuid |
