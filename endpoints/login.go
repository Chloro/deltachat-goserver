package endpoints

import(
	"bytes"
	"deltachat/models"
	"deltachat/keys"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
	"unicode/utf8"
)

func Login(request *http.Request, response http.ResponseWriter, s *mgo.Session) {

	var unauthenticatedUser models.UnauthenticatedUser
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&unauthenticatedUser)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	validateLogin(unauthenticatedUser)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	session := s.Copy()
	defer session.Close()
	collection := session.DB("deltachat").C("users")
	var existentUser models.UnauthenticatedUser
	err = collection.Find(bson.M{"username": unauthenticatedUser.Username}).One(&existentUser)

	if existentUser.Username == "" {
		http.Error(response, "That user does not exist.", http.StatusNotFound)
		return
	} else {

		err := bcrypt.CompareHashAndPassword(existentUser.HashedPassword, []byte(unauthenticatedUser.Password))
		if err != nil {
			http.Error(response, "Invalid password", http.StatusUnauthorized)
			return
		}

		existentUser.HashedPassword = nil
		signer := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["iss"] = "deltachat-server"
		claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
		claims["UserInfo"] = existentUser
		signer.Claims = claims

		tokenString, err := signer.SignedString(keys.SignKey)
		if err != nil {
			http.Error(response, "Error signing the token.", http.StatusUnauthorized)
			return
		}

		authenticatedUser := models.AuthenticatedUser{Username: existentUser.Username, Token: tokenString, Email: existentUser.Email}
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.WriteHeader(201)
		json.NewEncoder(response).Encode(authenticatedUser)
	}
}

func validateLogin(unauthenticatedUser models.UnauthenticatedUser) error {
	var errorBuffer bytes.Buffer

	if utf8.RuneCountInString(unauthenticatedUser.Username) == 0 {
		errorBuffer.WriteString("Username is a required field. \n")
	} else if utf8.RuneCountInString(unauthenticatedUser.Username) < 6 {
		errorBuffer.WriteString("Username should be at least 6 characters long. \n")
	} else if utf8.RuneCountInString(unauthenticatedUser.Username) > 32 {
		errorBuffer.WriteString("Username should be at most 32 characters long. \n")
	}

	if utf8.RuneCountInString(unauthenticatedUser.Password) == 0 {
		errorBuffer.WriteString("Password is a required field. \n")
	} else if utf8.RuneCountInString(unauthenticatedUser.Password) < 6 {
		errorBuffer.WriteString("Password should be at least 6 characters long. \n")
	} else if utf8.RuneCountInString(unauthenticatedUser.Password) > 64 {
		errorBuffer.WriteString("Password should be at most 64 characters long. \n")
	}

	if len(errorBuffer.String()) > 0 {
		return errors.New(errorBuffer.String())
	} else {
		return nil
	}
}
