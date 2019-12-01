package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb_ "gitlab.tesonero-computers.ru/ibis/aaa/endpoint"
	"google.golang.org/grpc"
)

type EdgeClient struct{}

var ServerUuid string

var addr = flag.String("addr1", "localhost:8448", "the server address to connect")
var addrUser = flag.String("addr3", "localhost:8447", "the user address to connect")

func AuthBackendSrv(client pb_.EdgeClient, account, ipAddress, port int32) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	req := &pb_.AuthRequest{
		Params: &pb_.EndpointParams{
			Account:  int32(account),
			AuthType: int32(2),
			// Timestamp:     m.Timestamp,
			// Rnd:           m.Rnd[:],
			// Hash:          m.Hash[:],
			// Cost:          int32(m.Cost),
			// Chid:          int32(m.ChID),
			// Hellointerval: int32(m.HelloInterval),
			RemoteAddr: int32(ipAddress),
			RemotePort: int32(port),
		},
	}

	resp, err := client.AuthBackendEndpoint(ctx, req)
	if err != nil {
		fmt.Println(err)
	}

	ServerUuid = resp.Uuid;

	fmt.Println(resp.Params)
	fmt.Println("Server Uuid:" + resp.Uuid)

}

func (s *EdgeClient) AuthEndpoint(ctx context.Context, in *pb_.AuthRequest) (*pb_.AuthResponse, error) {
	resp := &pb_.AuthResponse{}
	return resp, nil
}


//func (s *EdgeClient) SessionEndPoint(ctx context.Context, in *pb_.SessionRequest) (*pb_.SessionResponse, error) {
//	resp := &pb_.SessionResponse{}
//	return resp, nil
//}


func BackendSessActions(client pb_.EdgeClient, sessType bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	sessOpen  := false
	sessClose := false

	if sessType {
		sessOpen = true
	} else {
		sessClose = true
	}


	req := &pb_.SessionRequest{ Uuid: ServerUuid, Open: sessOpen, Close: sessClose}

	resp, err := client.SessionEndPoint(ctx, req)
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println("Server Session: OK")
	fmt.Println(resp)

}



// ########  USER OPERATION

func authUser(client pb_.UserClient, login string, password string, uid string) {
	fmt.Println("AuthUserEndPoint ----- ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.AuthUserEndPoint(ctx, &pb_.AuthUserRequest{SrvUuid: uid, Login: login, Password: password})
	if err != nil {
		log.Fatalf("client.AuthUserEndPoint(_) = _, %v: ", err)
	}


	fmt.Println("User Uuid: ", resp.UserUuid)
	fmt.Println("Response:  ", resp)
	fmt.Println("Message:   ", resp.Message)
	fmt.Println("Status:    ", resp.Open)
	//fmt.Println("Response Uid: ", resp.Uid)
}

func registerUser(client pb_.UserClient, login string, password string, uid string) {

	fmt.Println("RegistrationUserEndPoint ---- ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.RegistrationUserEndPoint(ctx, &pb_.AuthUserRequest{SrvUuid: uid, Login: login, Password: password})
	if err != nil {
		log.Fatalf("client.AuthUserEndPoint(_) = _, %v: ", err)
	}
	fmt.Println("Response: ", resp)
	fmt.Println("Message:  ", resp.Message)
	fmt.Println("Status:   ", resp.Open)
	//fmt.Println("Response Uid: ", resp.Uid)
}




func main() {


	//#######################
	// ---- BACKEND ACTION

	fmt.Println("I am Backend Client Start!")

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb_.NewEdgeClient(conn)

	var account   int32 = 1    // Важный параметр - account
	var ipAddress int32 = 444  // не используется
	var port      int32 = 777  // не используется

	AuthBackendSrv(client, account, ipAddress, port)  // авторизация сервера


	BackendSessActions(client, true)   //  открываем сессию


	// BackendSessActions(client, false)  //  закрываем сессию


    //#####################
	//---- USER  ACTION

	fmt.Println("I am User Client Start!")

	userConn, er := grpc.Dial(*addrUser, grpc.WithInsecure())
	if er != nil {
		log.Fatalf("did not connect: %v", er)
	}
	defer userConn.Close()


	userClient := pb_.NewUserClient(userConn)
	srvUid     := ServerUuid

	newLogin    := "Freedman Pol"
	newPassword := "1234"

	// регистрация нового пользователя
	// registerUser(userClient, newLogin, newPassword, srvUid)


	// Авторизация пользователя
	uLogin     := newLogin
	uPassword  := newPassword
	authUser(userClient, uLogin, uPassword, srvUid)


}
