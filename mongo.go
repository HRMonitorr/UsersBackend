package UsersBackend

import (
	"context"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func MongoCreateConnection(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: MongoString,
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func InsertDataEmployee(MongoConn *mongo.Database, colname string, emp Employee) (InsertedID interface{}) {
	return pasproj.InsertOneDoc(MongoConn, colname, emp)
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

func InsertOtp(MongoConn *mongo.Database, colname string, otp OTP) (InsertedID interface{}) {
	return pasproj.InsertOneDoc(MongoConn, colname, otp)
}

func GetOtp(mongoconn *mongo.Database, colname, otp string) (dest OTP) {
	filter := bson.M{"otp-code": otp}
	dest = atdb.GetOneDoc[OTP](mongoconn, colname, filter)
	return
}

func GetOtpExists(mongoconn *mongo.Database, colname, otp string) (exists bool) {
	var dest OTP
	dest = GetOtp(mongoconn, colname, otp)
	if dest.OTPCode == "" {
		return false
	}
	return true
}

func DeleteOTP(mcon *mongo.Database, colname, otp string) (deletedid *mongo.DeleteResult) {
	filter := bson.M{"otp-code": otp}
	deletedid = atdb.DeleteOneDoc(mcon, colname, filter)
	return
}

func InsertUserdata(MongoConn *mongo.Database, user Users) (InsertedID interface{}) {
	return pasproj.InsertOneDoc(MongoConn, "user", user)
}

func GetOneUser(MongoConn *mongo.Database, colname string, userdata Users) Users {
	filter := bson.M{"username": userdata.Username}
	data := atdb.GetOneDoc[Users](MongoConn, colname, filter)
	return data
}

func GetCommitwithusername(MongoConn *mongo.Database, colname, username string) (dest []structure.Commits) {
	filter := bson.M{"author": username}
	dest = atdb.GetAllDocByFilter[[]structure.Commits](MongoConn, colname, filter)
	return
}

func InsertWageData(MongoConn *mongo.Database, wage WageCalc) (InsertedID interface{}) {
	return pasproj.InsertOneDoc(MongoConn, "wage", wage)
}

func GetWgebyMonth(MongoConn *mongo.Database, month, name string) bool {
	filter := bson.M{"month": month, "employeeName": name}
	data := atdb.GetOneDoc[WageCalc](MongoConn, "wage", filter)
	if data.Month == time.Now().String() {
		return false
	}
	return true
}

func GetWgeAll(MongoConn *mongo.Database) []WageCalc {
	data := atdb.GetAllDoc[[]WageCalc](MongoConn, "wage")
	return data
}
