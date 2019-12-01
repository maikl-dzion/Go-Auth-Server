package backend

import (
	"context"
	"fmt"

	h_ "gitlab.tesonero-computers.ru/ibis/aaa/internal/handler"
	pb_ "gitlab.tesonero-computers.ru/ibis/aaa/endpoint"
	_ "gitlab.tesonero-computers.ru/ibis/aaa/internal/data"
	"gitlab.tesonero-computers.ru/ibis/aaa/internal/model"
)

type UserServer  struct{}

type BackendSrv  struct{}

var  BackEndUuid string


func (s *UserServer) AuthUserEndPoint(ctx context.Context,in *pb_.AuthUserRequest) (*pb_.AuthUserResponse, error) {

	srvUuid  := in.SrvUuid
	userUuid := ""
	message  := ""
	open  := false
	close := false

	// fmt.Println(BackEndUuid)

	resp := &pb_.AuthUserResponse {
		SrvUuid : srvUuid,
		UserUuid: userUuid,
		Login   : in.Login,
		Open    : open,
		Close   : close,
		Message : message,
	}

	// check server uuid
	var checkUuid = h_.CheckSrvUuid(in.SrvUuid, BackEndUuid)
	if !checkUuid {
		resp.Message = "Error srv uuid "
		fmt.Println(resp.Message)
		return resp, nil
	}


	//fmt.Println(" in server Uuid ", in.SrvUuid)
	//fmt.Println(" backend Uuid ", BackEndUuid)
	//fmt.Println(" checkUuid ", checkUuid)

	// auth user
	open = h_.AuthUserGrpc(in.Login, in.Password)
	if !open {
		resp.Message = "User authentication failed error"
		fmt.Println(resp)
		return resp, nil
	}


	userUuid = h_.CreateUuid()  //  create user uuid

	resp.Open     = open
	resp.UserUuid = userUuid
	resp.Message  = "User authentication : Ok"


	fmt.Println(resp.Message)

	//fmt.Println("srv uuid", srvUuid)
	//fmt.Println("user uuid", userUuid)

	//fmt.Println(in)
	//fmt.Println(srvUuid)
	//fmt.Println(resp.Open)
	//fmt.Println(resp.Message)

	return resp, nil

}



func (s *UserServer) RegistrationUserEndPoint(ctx context.Context, in *pb_.AuthUserRequest) (*pb_.AuthUserResponse, error) {


	// resp := &pb_.AuthUserResponse{}

	resp := &pb_.AuthUserResponse {
		SrvUuid  : in.SrvUuid,
		UserUuid : "",
		Login    : in.Login,
		Open     : false,
		Close    : false,
	}

	// check server uuid
	checkUuid := h_.CheckSrvUuid(in.SrvUuid, BackEndUuid)
	if !checkUuid {
		resp.Message = "Error srv uuid"
		return resp, nil
	}


	// user register
	state, _ := h_.RegisterUserGrpc(in.Login, in.Password)

	switch state {
	case 1 :
		resp.Message = "login error"
		return resp, nil

		//case 2 :
		//	resp.Message = "login error"
		//	return resp, nil

	case 3 :
		resp.Message = "db save error"
		return resp, nil

	default :
		resp.Message = "User register : Ok"
	}


	// create user uuid
	userUuid     := h_.CreateUuid()
	resp.Open     = true
	resp.UserUuid = userUuid
	// resp.Message  = "User register Ok"
	fmt.Println(resp.Message)

	return resp, nil
}



func (s *BackendSrv) AuthBackendEndpoint(ctx context.Context,in *pb_.AuthRequest) (*pb_.AuthResponse, error) {

	BackEndUuid = ""

	resp := &pb_.AuthResponse {
		Params: in.Params,
		Uuid: BackEndUuid,
	}

	// fmt.Println(in.Params.Account)
	// acc := in.Params.Account
	authType := in.Params.AuthType
	srvAccount := in.Params.Account

	//if authType == 5 {
	//
	//	sess := "de83a4b7-7d4f-11e9-991d-00ff7554c695"
	//	_ , err := h_.CloseBackendSession(sess)
	//	if err != nil {
	//		fmt.Println("BackendServer Session Close : ERROR ")
	//		return resp, nil
	//	}
	//	fmt.Println("BackendServer Session Close : OK ")
	//	return resp, nil
	//}


	if authType != 2 || srvAccount != model.BackendAccount {
		fmt.Println("Backend Auth : ERROR ")
		return resp, nil
	}


	BackEndUuid = h_.CreateUuid()

    //_ , err := h_.SaveBackendSession(BackEndUuid)
    //if err != nil {
	//	fmt.Println("BackendServer Session Save : ERROR ")
	//	return resp, nil
	//}

	resp.Uuid = BackEndUuid

	fmt.Println("Backend Auth : Ok ")

	return resp, nil

}


func (s *BackendSrv) SessionEndPoint(ctx context.Context,in *pb_.SessionRequest) (*pb_.SessionResponse, error) {

	resp := &pb_.SessionResponse {}

	// sess := "de83a4b7-7d4f-11e9-991d-00ff7554c695"
	// sessOpen := in.Open

	srvUuid := in.Uuid
	messageAction := "Backend Session "

	if in.Open {
		messageAction = messageAction + " Save"
		_ , err := h_.SaveBackendSession(srvUuid)
		if err != nil {
			fmt.Println(messageAction + ": ERROR ")
			return resp, nil
		}

		fmt.Println(messageAction + " : OK ")

	}


	if in.Close {
		messageAction = messageAction + " Close"
		_ , err := h_.CloseBackendSession(srvUuid)
		if err != nil {
			fmt.Println(messageAction + " : ERROR ")
			return resp, nil
		}
		fmt.Println(messageAction + ": OK ")

	}


	return resp, nil

}


func (s *BackendSrv) AuthEndpoint(ctx context.Context,in *pb_.AuthRequest) (*pb_.AuthResponse, error) {

	resp := &pb_.AuthResponse {}

	return resp, nil

}

