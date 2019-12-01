package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uid_ "github.com/satori/go.uuid"
	"gitlab.tesonero-computers.ru/ibis/aaa/internal/data"
	"gitlab.tesonero-computers.ru/ibis/aaa/internal/model"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	users := data.GetUsers()

	json.NewEncoder(w).Encode(users)

}


func GetUserById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	userId, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	user, err := data.GetUserById(userId)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("db select row error... %v", err)
		return
	}

	json.NewEncoder(w).Encode(user)

}


func AuthUser(w http.ResponseWriter, r *http.Request) {

	// defer r.Body.Close()

	// decoder := json.NewDecoder(r.Body)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	authUser := false

	u := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("json decode error... %v", err)
		return
	}

	// fmt.Println(u)

	user, err := data.FindLogin(u.Login);

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("db select error... %v", err)
		json.NewEncoder(w).Encode(struct{result bool }{ result:authUser})
		return
	}


	// fmt.Println(u.Password)
	// fmt.Println(user.Password)

	authUser = passwdCompare(u.Password, user.Password)


	fmt.Println(authUser)

	json.NewEncoder(w).Encode(struct{result bool}{ result:authUser})
}


func RegisterUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	u := &model.User{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("json decode error... %v", err)
		return
	}


	//fmt.Println(u)


	id, err := data.RegisterUser(u)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("db save error... %v", err)
		return
	}

	u.Id = id

	json.NewEncoder(w).Encode(u)

}



func DeleteUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := mux.Vars(r)

	err := data.DeleteUser(params["login"])
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("db delete user error... %v", err)
		return
	}

	// fmt.Println("delete user OK")

	json.NewEncoder(w).Encode(struct{result bool}{ result:true})
}



func ChangeUserPasswd(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ch := model.ChangePasswd{}

	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("json decode error... %v", err)
		return
	}


    us, err := data.ChangeUserPasswd(ch)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(res)
	//fmt.Println(ch)
	json.NewEncoder(w).Encode(us)

}


func passwdCompare(reqPasswd string, realPasswd string) (bool){

	auth := false

	if(reqPasswd == realPasswd) {
		auth = true;
	}

	return auth

}



func CheckBody(w http.ResponseWriter, r *http.Request) (bool) {
	ret := true;
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		ret = false;
	}

	return ret
}



// ############################
//  ---------   GRPC  --------


func AuthUserGrpc(login string, password string) (bool){

    // var message string = "Ok"
    var auth bool = false

	u, err := data.FindLogin(login);
	// fmt.Println(u)
	if err != nil {
		log.Printf("db select error... %v", err)
		return auth
	}

	if password == u.Password {
		auth = true
	}

	return auth
}



func RegisterUserGrpc(login, password string) (int, int) {

	u := &model.User{}
	u.Password = password
	u.Login    = login

	_, err := data.FindLogin(login)
	if err == nil {
		log.Printf("login error... %v", err)
		return 1, 0
	}


	//if user.Login == login {
	//	log.Printf("login error... ")
	//	return 2, 0
	//}

	insertId, err := data.RegisterUser(u)
	if err != nil {
		log.Printf("db save error... %v", err)
		return 3, 0
	}

	return 0, insertId

}


func CheckSrvUuid(srvUuid, realSrvUuid string) (bool){

	result := false

	// realSrvUuid := "645378383"

	if srvUuid == "" {
		return result
	}

	_ , err := data.GetAccSession(srvUuid)
	if err != nil {
		fmt.Println("Error : Session not ", err)
		return result
	}


	//if resp.ClosedAt != "" {
	//	fmt.Println("Error : Session close ", resp.ClosedAt)
	//	return result
	//}

	//if srvUuid == "" || realSrvUuid == "" {
	//	return result
	//}
	//
	//if srvUuid == realSrvUuid {
	//	result = true
	//}

	return true
}


func CreateUuid() (string) {

	newUuid := uid_.NewV1().String()

    return newUuid
}



func SaveBackendSession(uuid string) (int, error){
	backendAccount := model.BackendAccount
	dt := time.Now()
	currDate := dt.Format("01-02-2006 15:04:05")
    id, err  := data.IbisSessionAdd(uuid, currDate, backendAccount)
    return id, err
}



func CloseBackendSession(uuid string) (int, error){
	backendAccount := model.BackendAccount
	dt := time.Now()
	currDate := dt.Format("01-02-2006 15:04:05")
	id, err  := data.IbisSessionClose(uuid, currDate, backendAccount)
	return id, err
}


//func GetBackendSession(uuid string) (int, error){
//	backendAccount := model.BackendAccount
//	dt := time.Now()
//	currDate := dt.Format("01-02-2006 15:04:05")
//	id, err  := data.IbisSessionAdd(uuid, currDate, backendAccount)
//	return id, err
//}