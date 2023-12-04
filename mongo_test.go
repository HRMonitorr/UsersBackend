package UsersBackend

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertOtp(t *testing.T) {
	table := MongoCreateConnection("mongodb+srv://rofinafiin:aXz4RdVqUVIQcqa1@rofinafiinsdata.9fyvx4r.mongodb.net", "HRMApp")
	data := OTP{
		Username: "rofi",
		DateOTP:  time.Now(),
		OTPCode:  CreateOTP(),
	}
	ins := InsertOtp(table, "otp", data)
	fmt.Printf("result : %s", ins)
}

func TestInsertUserdata(t *testing.T) {
	table := MongoCreateConnection("mongodb+srv://rofinafiin:aXz4RdVqUVIQcqa1@rofinafiinsdata.9fyvx4r.mongodb.net", "HRMApp")
	data := Users{
		Username: "rofi",
		Password: "password",
		PhoneNum: "6285156007137",
		Role:     "admin",
	}
	ins := InsertUserdata(table, data)
	fmt.Printf("result : %s", ins)
}
