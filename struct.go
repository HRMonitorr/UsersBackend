package UsersBackend

import (
	"time"
)

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

type Users struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	PhoneNum string `json:"phone-num" bson:"phoneNum"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type Logindata struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type ResponseEmployeeBanyak struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Employee `json:"data"`
}

type Employee struct {
	EmployeeId string   `json:"employeeid" bson:"employeeid,omitempty"`
	Name       string   `json:"name" bson:"name,omitempty"`
	Username   string   `json:"username" bson:"username"`
	Email      string   `json:"email" bson:"email,omitempty"`
	Phone      string   `json:"phone" bson:"phone,omitempty"`
	Division   Division `json:"division" bson:"division,omitempty"`
	Salary     Salary   `json:"salary" bson:"salary"`
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

type RequestEmployee struct {
	EmployeeId string `json:"employeeid"`
}

type Responses struct {
	Response string `bson:"response" json:"response"`
}

type OTP struct {
	Username string    `json:"username" bson:"username"`
	Role     string    `bson:"role" json:"role"`
	DateOTP  time.Time `json:"date-otp" bson:"date-otp"`
	OTPCode  string    `bson:"otp-code" json:"otp-code"`
}

type OnlyOTP struct {
	OTPCode string `bson:"otp-code" json:"otp-code"`
}

type WageCalc struct {
	EmployeeName    string  `bson:"employeeName" json:"employeeName"`
	JumlahCommit    float64 `bson:"jumlahCommit" json:"jumlahCommit"`
	BasicSalary     float64 `bson:"basicSalary" json:"basicSalary"`
	HonorDivision   float64 `bson:"honorDivision" json:"honorDivision"`
	InsentifCommits float64 `json:"insentifCommits" bson:"insentifCommits"`
	GrossSalary     float64 `bson:"jumlahKotor" json:"jumlahKotor"`
	Tax             float64 `json:"tax" bson:"tax"`
	NettSalary      float64 `bson:"nettSalary" json:"nettSalary"`
	Month           string  `json:"month" bson:"month"`
}
