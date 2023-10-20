package UsersBackend

import pasproj "github.com/HRMonitorr/PasetoprojectBackend"

type ResponseBack struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type ResponseEmployee struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Employee `json:"data"`
}

type Employee struct {
	EmployeeId string       `json:"employeeid" bson:"employeeid"`
	Name       string       `json:"name" bson:"name"`
	Email      string       `json:"email" bson:"email"`
	Phone      string       `json:"phone" bson:"phone"`
	Division   Division     `json:"division" bson:"division"`
	Account    pasproj.User `json:"account" bson:"account"`
}

type Division struct {
	DivId   int    `json:"divId" bson:"divId"`
	DivName string `json:"divName" bson:"divName"`
}

type Updated struct {
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}
