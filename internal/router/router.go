package router

import (
	"github.com/gorilla/mux"

	"gitlab.tesonero-computers.ru/ibis/aaa/internal/handler"
)

func RoutesInit() *mux.Router {
	mx := mux.NewRouter()

	mx.HandleFunc("/users", handler.GetUsers).Methods("GET")
	mx.HandleFunc("/users/{id}", handler.GetUserById).Methods("GET")
	mx.HandleFunc("/users/register", handler.RegisterUser).Methods("POST")
	mx.HandleFunc("/users/auth", handler.AuthUser).Methods("POST")
	mx.HandleFunc("/users/{login}", handler.DeleteUser).Methods("DELETE")
	mx.HandleFunc("/users/change_passwd", handler.ChangeUserPasswd).Methods("PUT")

	mx.HandleFunc("/app/auth", handler.CreateTokenEndpoint).Methods("POST")
	mx.HandleFunc("/app/check", handler.ProtectedEndpoint).Methods("GET")
	//mx.HandleFunc("/users/add", addUser).Methods("GET")

	return mx
}
