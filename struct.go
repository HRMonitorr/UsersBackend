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

type ResponseEmployeeBanyak struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Employee `json:"data"`
}

type Employee struct {
	EmployeeId string       `json:"employeeid" bson:"employeeid,omitempty"`
	Name       string       `json:"name" bson:"name,omitempty"`
	Email      string       `json:"email" bson:"email,omitempty"`
	Phone      string       `json:"phone" bson:"phone,omitempty"`
	Division   Division     `json:"division" bson:"division,omitempty"`
	Account    pasproj.User `json:"account" bson:"account,omitempty"`
	Salary     Salary       `json:"salary" bson:"salary"`
}

type Division struct {
	DivId   int    `json:"divId" bson:"divId"`
	DivName string `json:"divName" bson:"divName"`
}

type Updated struct {
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type Salary struct {
	BasicSalary   int `bson:"basic-salary" json:"basic-salary"`
	HonorDivision int `bson:"honor-division" json:"honor-division"`
}

type Cred struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ReqUsers struct {
	Username string `json:"username"`
}
