package main

import (
	"deltachat/keys"
	"fmt"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

func main() {
	keys.InitKeys()
	//////////////Start up database/////////////
	fmt.Println("Connecting to the database...")
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Errorf("Failed to connect to MongoDB: %s", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	fmt.Println("Connected to mongo!")
	///////////////////////////////////////////
	router := mux.NewRouter()
	registerEndpoints(router, session)
	log.Fatal(http.ListenAndServe(":9001", router))
}
