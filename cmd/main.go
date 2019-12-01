package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"sync"

	pb_ "gitlab.tesonero-computers.ru/ibis/aaa/endpoint"
	_ "gitlab.tesonero-computers.ru/ibis/aaa/internal/data"

	ibis_h "gitlab.tesonero-computers.ru/ibis/aaa/internal/ibis"
	ibispb "gitlab.tesonero-computers.ru/ibis/authproto/pkg/endpoint"


	back_h "gitlab.tesonero-computers.ru/ibis/aaa/internal/backend"

	// "gitlab.tesonero-computers.ru/ibis/aaa/internal/router"

)


var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)



func setEnvValues() {

	//os.Setenv("DB_ADDRESS", "192.168.3.23")
	//os.Setenv("DB_PORT", "5432")
	//os.Setenv("DB_ROLE", "ibis")
	//os.Setenv("DB_PASSWD", "ibis")
	//os.Setenv("DB_DATABASE", "ibis")


	os.Setenv("AAA_IBIS_PORT", ":8443")
	os.Setenv("AAA_USER_PORT", ":8447")
	os.Setenv("AAA_BACKEND_PORT", ":8448")
}



func main() {

	flag.Parse()

	setEnvValues()

	ibisPort := os.Getenv("AAA_IBIS_PORT")
	userPort := os.Getenv("AAA_USER_PORT")
	backendPort := os.Getenv("AAA_BACKEND_PORT")


	var wg sync.WaitGroup

	// ######################
	// --- USER SERVER START

	// listenUser, err := net.Listen("tcp", ":8445")
	listenUser, err := net.Listen("tcp", userPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb_.RegisterUserServer(srv, &back_h.UserServer{})

	fmt.Println("UserServer starting on port ..." + userPort)

	wg.Add(1)
	go func () {
		defer wg.Done()
		if err := srv.Serve(listenUser); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()



    // #########################
	// --- BACKEND SERVER START

	listenBackEnd, err := net.Listen("tcp", backendPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	backSrv := grpc.NewServer()
	pb_.RegisterEdgeServer(backSrv, &back_h.BackendSrv{})


	fmt.Println("BackendSrv starting on port ..." + backendPort)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := backSrv.Serve(listenBackEnd); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()



	// #########################
	// --- IBIS SERVER START

	ibisSrv := grpc.NewServer()
	// defer srv.GracefulStop()

	ibispb.RegisterAAAServer(ibisSrv, &ibis_h.Server{})

	ibisListener, err := net.Listen("tcp", ibisPort)
	if err != nil {
		// log.Fatal().Err(err).Msg("")
		log.Fatal(err)
	}

	fmt.Println("ibisServer starting on port ..." + ibisPort)


	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = ibisSrv.Serve(ibisListener); err != nil {
			log.Fatal(err)
		}

    }()



    /////////////////////////////
	wg.Wait()


	//mx := router.RoutesInit()
	//fmt.Println(" ...Server init Ok... ")
	//log.Fatal(http.ListenAndServe(":8989", mx))

}

