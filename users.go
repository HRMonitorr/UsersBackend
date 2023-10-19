package UsersBackend

import (
	"encoding/json"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/whatsauth/watoken"
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

func GetDataUserForAdmin(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(pasproj.ResponseDataUser)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	cihuy := new(pasproj.Response)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	}

	cekadmin := IsAdmin(tokenlogin, PublicKey)
	if cekadmin != true {
		req.Status = false
		req.Message = "IHHH Kamu bukan admin"
	}

	err := json.NewDecoder(r.Body).Decode(&cihuy)
	if err != nil {
		req.Status = false
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		checktoken := watoken.DecodeGetId(os.Getenv(PublicKey), cihuy.Token)
		compared := pasproj.CompareUsername(conn, colname, checktoken)
		if compared != true {
			req.Status = false
			req.Message = "Data User tidak ada"
		} else {
			datauser := pasproj.GetAllUser(conn, colname)
			req.Status = true
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return pasproj.ReturnStringStruct(req)
}
