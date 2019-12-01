package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type AppAuth struct {
	Appname string `json:"appname"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {

	var app AppAuth

	_ = json.NewDecoder(req.Body).Decode(&app)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"appname": app.Appname,
		"password": app.Password,
	})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}


func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {

	params := req.URL.Query()

	// fmt.Println(params)

	token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var app AppAuth
		mapstructure.Decode(claims, &app)
		json.NewEncoder(w).Encode(app)

	} else {

		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})

	}
}


//func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		authorizationHeader := req.Header.Get("authorization")
//		if authorizationHeader != "" {
//			bearerToken := strings.Split(authorizationHeader, " ")
//			if len(bearerToken) == 2 {
//				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
//					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//						return nil, fmt.Errorf("There was an error")
//					}
//					return []byte("secret"), nil
//				})
//				if error != nil {
//					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
//					return
//				}
//				if token.Valid {
//					context.Set(req, "decoded", token.Claims)
//					next(w, req)
//				} else {
//					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
//				}
//			}
//		} else {
//			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
//		}
//	})
//}
//
//
//func TestEndpoint(w http.ResponseWriter, req *http.Request) {
//	decoded := context.Get(req, "decoded")
//	var user User
//	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
//	json.NewEncoder(w).Encode(user)
//}
