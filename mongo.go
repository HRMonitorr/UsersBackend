package UsersBackend

import (
	"context"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDataEmployee(MongoConn *mongo.Database, colname string, emp Employee) (InsertedID interface{}) {
	req := new(Employee)
	req.EmployeeId = emp.EmployeeId
	req.Name = emp.Name
	req.Email = emp.Email
	req.Phone = emp.Phone
	req.Division = emp.Division
	req.Account = emp.Account
	return pasproj.InsertOneDoc(MongoConn, colname, req)
}

func GetAllEmployeeData(Mongoconn *mongo.Database, colname string) []Employee {
	data := atdb.GetAllDoc[[]Employee](Mongoconn, colname)
	return data
}

func DeleteUser(Mongoconn *mongo.Database, colname, username string) (deleted interface{}, err error) {
	filter := bson.M{"username": username}
	data := atdb.DeleteOneDoc(Mongoconn, colname, filter)
	return data, err
}

func UpdateEmployee(Mongoconn *mongo.Database, ctx context.Context, emp Employee) (UpdateId interface{}, err error) {
	filter := bson.D{{"employeeid", emp.EmployeeId}}
	res, err := Mongoconn.Collection("employee").ReplaceOne(ctx, filter, emp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdatePassword(mongoconn *mongo.Database, user pasproj.User) (Updatedid interface{}) {
	filter := bson.D{{"username", user.Username}}
	pass, _ := pasproj.HashPass(user.Password)
	update := bson.D{{"$Set", bson.D{
		{"password", pass},
	}}}
	res, err := mongoconn.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "gagal update data"
	}
	return res
}

func DeleteEmployeeData(mongoconn *mongo.Database, colname, EmpId string) (deletedid interface{}, err error) {
	filter := bson.M{"employeeid": EmpId}
	data := atdb.DeleteOneDoc(mongoconn, colname, filter)
	return data, err
}

func GetOneEmployeeData(mongoconn *mongo.Database, colname, Empid string) (dest Employee) {
	filter := bson.M{"employeeid": Empid}
	dest = atdb.GetOneDoc[Employee](mongoconn, colname, filter)
	return
}
