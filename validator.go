package UsersBackend

import (
	"fmt"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"math/rand"
)

func IsAdmin(Tokenstr, PublicKey string) bool {
	role, err := pasproj.DecodeGetRole(PublicKey, Tokenstr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "admin" {
		return false
	}
	return true
}

func IsHR(TokenStr, Publickey string) bool {
	role, err := pasproj.DecodeGetRole(Publickey, TokenStr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "HR" {
		return false
	}
	return true
}

func CreateOTP() string {
	return RandStringBytes(6)
}

// generateRandomString generates a random string of a specified length using the given characters.
const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXZ1238849103748102"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
