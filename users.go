package UsersBackend

import (
	"context"
	"encoding/json"
	"fmt"
	pasproj "github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"github.com/aiteung/atapi"
	"github.com/gofiber/fiber/v2"
	"github.com/whatsauth/wa"
	"net/http"
	"os"
	"time"
)

// reg User
func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(pasproj.Credential)
	userdata := new(Users)
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
		data := pasproj.User{
			Username: userdata.Username,
			Password: hash,
			PhoneNum: userdata.PhoneNum,
			Role:     userdata.Role,
		}
		pasproj.InsertUserdata(conn, data)
		resp.Message = "Berhasil Input data"
	}
	response := pasproj.ReturnStringStruct(resp)
	return response

}

// log User
func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp pasproj.Credential
	mconn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	var datauser pasproj.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if pasproj.PasswordValidator(mconn, Colname, datauser) {
			datarole := pasproj.GetOneUser(mconn, "user", pasproj.User{Username: datauser.Username})
			tokenstring, err := pasproj.EncodeWithRole(datarole.Role, datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

func LoginOTP(MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp pasproj.Credential
	mconn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	var datauser Logindata
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		datarole := GetOneUser(mconn, "user", Users{Username: datauser.Username})
		if datarole.Username == "" {
			resp.Message = "Data User tidak ditemukan nih bestie"
		}
		if pasproj.PasswordValidator(mconn, Colname, pasproj.User{
			Username: datauser.Username,
			Password: datauser.Password,
			Role:     datarole.Role,
		}) {
			data := OTP{
				Username: datauser.Username,
				Role:     datarole.Role,
				DateOTP:  time.Now(),
				OTPCode:  CreateOTP(),
			}
			InsertOtp(mconn, "otp", data)
			dt := &wa.TextMessage{
				To:       datarole.PhoneNum,
				IsGroup:  false,
				Messages: fmt.Sprintf("Hai hai kak Ini OTP kakak %s", data.OTPCode),
			}
			res, _ := atapi.PostStructWithToken[Responses]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
			resp.Status = true
			resp.Message = res.Response
			resp.Token = data.OTPCode
		} else {
			resp.Message = "Password Salah"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

func LoginAfterOTP(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp pasproj.Credential
	mconn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	var datauser OnlyOTP
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if GetOtpExists(mconn, "otp", datauser.OTPCode) {
			dataOTP := GetOtp(mconn, "otp", datauser.OTPCode)
			datarole := pasproj.GetOneUser(mconn, "user", pasproj.User{Username: dataOTP.Username})
			tokenstring, err := pasproj.EncodeWithRole(datarole.Role, dataOTP.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				DeleteOTP(mconn, "otp", dataOTP.OTPCode)
				resp.Status = true
				resp.Message = "Berhasil Login, Selamat datangg"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "OTP Salah"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

// Get Data User
func GetDataUserForAdmin(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(pasproj.ResponseDataUser)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		cekadmin := IsAdmin(tokenlogin, PublicKey)
		if cekadmin != true {
			req.Status = false
			req.Message = "IHHH Kamu bukan admin"
		}
		checktoken, err := pasproj.DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "tidak ada data username : " + tokenlogin
		}
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

// Reset Password
func ResetPassword(MongoEnv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(pasproj.User)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			resp.Status = fiber.StatusInternalServerError
			resp.Message = "kamu bukan admin"
		} else {
			UpdatePassword(conn, pasproj.User{
				Username: req.Username,
				Password: req.Password,
			})
			resp.Status = fiber.StatusOK
			resp.Message = "Berhasil reset password"
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

// Delete User
func DeleteUserforAdmin(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(ReqUsers)
	conn := pasproj.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Message = "error parsing application/json: " + err.Error()
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				resp.Status = fiber.StatusInternalServerError
				resp.Message = "kamu bukan admin"
			} else {
				_, err := DeleteUser(conn, colname, req.Username)
				if err != nil {
					resp.Status = fiber.StatusBadRequest
					resp.Message = "gagal hapus data"
				}
				resp.Status = fiber.StatusOK
				resp.Message = "data berhasil dihapus"
			}
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

// Insert data
func InsertEmployee(MongoEnv, dbname, colname, publickey string, r *http.Request) string {
	resp := new(pasproj.Credential)
	req := new(Employee)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			checkHR := IsHR(tokenlogin, os.Getenv(publickey))
			if !checkHR {
				resp.Status = false
				resp.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
			}
		} else {
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp.Message = "error parsing application/json: " + err.Error()
			} else {
				InsertDataEmployee(conn, colname, Employee{
					EmployeeId: req.EmployeeId,
					Name:       req.Name,
					Email:      req.Email,
					Username:   req.Username,
					Phone:      req.Phone,
					Division: Division{
						DivId:   req.Division.DivId,
						DivName: req.Division.DivName,
					},
					Salary: Salary{
						BasicSalary:   req.Salary.BasicSalary,
						HonorDivision: req.Salary.HonorDivision,
					},
				})
				resp.Status = true
				resp.Message = "Berhasil Insert data"
			}
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

// Update data
func UpdateDataEmployees(MongoEnv, dbname, publickey string, r *http.Request) string {
	req := new(pasproj.Credential)
	resp := new(Employee)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				checkHR := IsHR(tokenlogin, os.Getenv(publickey))
				if !checkHR {
					req.Status = false
					req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
				}
			} else {
				conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
				UpdateEmployee(conn, context.Background(), Employee{
					EmployeeId: resp.EmployeeId,
					Name:       resp.Name,
					Email:      resp.Email,
					Username:   resp.Username,
					Phone:      resp.Phone,
					Division: Division{
						DivId:   resp.Division.DivId,
						DivName: resp.Division.DivName,
					},
					Salary: Salary{
						BasicSalary:   resp.Salary.BasicSalary,
						HonorDivision: resp.Salary.HonorDivision,
					},
				})
				req.Status = true
				req.Message = "Berhasil Update data"
			}
		}
	}
	return pasproj.ReturnStringStruct(req)
}

// Get One
func GetOneEmployee(PublicKey, MongoEnv, dbname string, r *http.Request) string {
	req := new(ResponseEmployee)
	resp := new(RequestEmployee)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Message = "error parsing application/json: " + err.Error()
	}
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			checkHR := IsHR(tokenlogin, os.Getenv(PublicKey))
			if !checkHR {
				req.Status = fiber.StatusBadRequest
				req.Message = "Anda tidak bisa Get data karena bukan HR atau admin"
			}
		} else {
			datauser := GetOneEmployeeData(conn, "employee", resp.EmployeeId)
			if datauser.EmployeeId == "" {
				req.Status = fiber.StatusBadRequest
				req.Message = "data User gagal diambil " + resp.EmployeeId
			}
			req.Status = fiber.StatusOK
			req.Message = "data User berhasil diambil " + resp.EmployeeId
			req.Data = datauser
		}
	}
	return pasproj.ReturnStringStruct(req)
}

// GetAll
func GetAllEmployee(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseEmployeeBanyak)
	conn := pasproj.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			checkHR := IsHR(tokenlogin, os.Getenv(PublicKey))
			if !checkHR {
				req.Status = fiber.StatusBadRequest
				req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
			}
		} else {
			datauser := GetAllEmployeeData(conn, colname)
			req.Status = fiber.StatusOK
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return pasproj.ReturnStringStruct(req)
}

// Delete Data
func DeleteEmployee(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(RequestEmployee)
	conn := pasproj.MongoCreateConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Message = "error parsing application/json: " + err.Error()
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				resp.Status = fiber.StatusInternalServerError
				resp.Message = "kamu bukan admin"
			} else {
				_, err := DeleteEmployeeData(conn, colname, req.EmployeeId)
				if err != nil {
					resp.Status = fiber.StatusBadRequest
					resp.Message = "gagal hapus data"
				}
				resp.Status = fiber.StatusOK
				resp.Message = "data berhasil dihapus"
			}
		}
	}
	return pasproj.ReturnStringStruct(resp)
}

func GetSalaryEmployee(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(structure.Creds)
	resp := new(RequestEmployee)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
			if !checkadmin {
				checkHR := IsHR(tokenlogin, os.Getenv(PublicKey))
				if !checkHR {
					req.Status = fiber.StatusBadRequest
					req.Message = "Anda tidak bisa Get data karena bukan HR atau admin"
				}
			} else {
				datauser := GetOneEmployeeData(conn, "employee", resp.EmployeeId)
				if datauser.EmployeeId == "" {
					req.Status = fiber.StatusBadRequest
					req.Message = "data user tidak ada"
				}
				dataCommit := GetCommitwithusername(conn, "commit", datauser.Username)
				if len(dataCommit) == 0 {
					if datauser.EmployeeId == "" {
						req.Status = fiber.StatusBadRequest
						req.Message = "data commit 0 "
					}
				}
				jumlahcommit := len(dataCommit)
				if jumlahcommit > 20 {
					jumlahcommit = 20
				}
				insentif := jumlahcommit * 20000
				tax := float64(datauser.Salary.BasicSalary+datauser.Salary.HonorDivision+insentif) * 0.15
				data := WageCalc{
					EmployeeName:    datauser.Name,
					JumlahCommit:    float64(len(dataCommit)),
					BasicSalary:     float64(datauser.Salary.BasicSalary),
					HonorDivision:   float64(datauser.Salary.HonorDivision),
					InsentifCommits: float64(insentif),
					GrossSalary:     float64(datauser.Salary.BasicSalary + datauser.Salary.HonorDivision + insentif),
					Tax:             tax,
					NettSalary:      float64(datauser.Salary.BasicSalary+datauser.Salary.HonorDivision+insentif) - tax,
					Month:           time.Now().Month().String(),
				}

				wagedata := GetWgebyMonth(conn, time.Now().Month().String(), datauser.Name)
				if wagedata {
					InsertWageData(conn, data)
				} else {
					req.Status = fiber.StatusBadRequest
					req.Message = "Data wage untuk bulan ini sudah ada " + resp.EmployeeId
				}

				req.Status = fiber.StatusOK
				req.Message = "data wage berhasil diinput " + datauser.Username + " " + resp.EmployeeId
				req.Data = data
			}
		}
	}
	return pasproj.ReturnStringStruct(req)
}

func GetSalaryAll(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(structure.Creds)
	conn := pasproj.MongoCreateConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			checkHR := IsHR(tokenlogin, os.Getenv(PublicKey))
			if !checkHR {
				req.Status = fiber.StatusBadRequest
				req.Message = "Anda tidak bisa Get data karena bukan HR atau admin"
			}
		} else {
			data := GetWgeAll(conn)
			req.Status = fiber.StatusOK
			req.Message = "data wage berhasil diinput"
			req.Data = data
		}
	}
	return pasproj.ReturnStringStruct(req)
}
