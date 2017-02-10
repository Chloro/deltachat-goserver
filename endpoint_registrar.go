package main

import (
	"deltachat/endpoints"
	//"deltachat/keys"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"net/http"
)

func registerEndpoints(router *mux.Router, session *mgo.Session) {
	router.HandleFunc("/api/login", func(response http.ResponseWriter, request *http.Request) {
		endpoints.Login(request, response, session)
	}).Methods("POST")

	router.HandleFunc("/api/register", func(response http.ResponseWriter, request *http.Request) {
		endpoints.RegisterUser(request, response, session)
	}).Methods("POST")

	router.HandleFunc("/api/users", func(response http.ResponseWriter, request *http.Request) {
		//ValidateTokenMiddleware(request, response)
		endpoints.GetAllUsers(request, response, session)
	}).Methods("GET")
}


//func ValidateTokenMiddleware(request *http.Request, response http.ResponseWriter) {
//
//	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error){
//		return keys.VerifyKey, nil
//	})
//
//	if err == nil {
//		if token.Valid{
//			next(response, r)
//		} else {
//			response.WriteHeader(http.StatusUnauthorized)
//			fmt.Fprint(response, "Token is not valid")
//		}
//	} else {
//		response.WriteHeader(http.StatusUnauthorized)
//		fmt.Fprint(response, "Unauthorised access to this resource")
//	}
//
//}