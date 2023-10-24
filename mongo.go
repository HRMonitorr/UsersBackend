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

func GetDataEmployee(MongoConn *mongo.Database, colname, empid string) Employee {
	filter := bson.M{"employeeid": empid}
	data := atdb.GetOneDoc[Employee](MongoConn, colname, filter)
	return data
}

func UpdateEmployee(Mongoenv, dbname string, ctx context.Context, emp Employee) (UpdateId interface{}) {
	conn := pasproj.MongoCreateConnection(Mongoenv, dbname)
	filter := bson.D{{"employeeid", emp.EmployeeId}}
	update := bson.D{{"$set", bson.D{
		{"phone", emp.Phone},
		{"email", emp.Email},
	}}}
	res, err := conn.Collection("employee").UpdateOne(ctx, filter, update)
	if err != nil {
		return "Gagal Update"
	}
	return res
}