package UsersBackend

import (
	"encoding/json"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"net/http"
	"os"
)

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(pasproj.Credential)
	userdata := new(pasproj.User)
	resp.Status = false
	conn := pasproj.MongoCreateConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := pasproj.HashPass(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		pasproj.InsertUserdata(conn, userdata.Username, userdata.Role, hash)
		resp.Message = "Berhasil Input data"
	}
	response := pasproj.ReturnStringStruct(resp)
	return response

}

func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp pasproj.Credential
	mconn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	var datauser pasproj.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if pasproj.PasswordValidator(mconn, Colname, datauser) {
			tokenstring, err := pasproj.EncodeWithRole(datauser.Role, datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}
