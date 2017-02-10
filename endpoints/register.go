package endpoints

import(
	"bytes"
	"golang.org/x/crypto/bcrypt"
	"deltachat/models"
	"encoding/json"
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"regexp"
	"unicode/utf8"
)

func RegisterUser(request *http.Request, response http.ResponseWriter, s *mgo.Session) {
	var unregisteredUser models.UnregisteredUser
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&unregisteredUser)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	err = validateRegistration(unregisteredUser)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	session := s.Copy()
	defer session.Close()
	collection := session.DB("deltachat").C("users")
	var existentUser models.UnauthenticatedUser
	err = collection.Find(bson.M{"username": unregisteredUser.Username}).One(&existentUser)

	if existentUser.Username != "" {
		http.Error(response, "That user already exists.", http.StatusConflict)
		return
	} else {
		unregisteredUser.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(unregisteredUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(response, "Error hashing password.", http.StatusInternalServerError)
			return
		}

		unregisteredUser.Password = ""

		err = collection.Insert(unregisteredUser)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		unregisteredUser.HashedPassword = nil
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.WriteHeader(201)
		json.NewEncoder(response).Encode(unregisteredUser)
	}
}

func validateRegistration(unregisteredUser models.UnregisteredUser) error {
	var errorBuffer bytes.Buffer

	if utf8.RuneCountInString(unregisteredUser.Username) == 0 {
		errorBuffer.WriteString("Username is a required field. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Username) < 6 {
		errorBuffer.WriteString("Username should be at least 6 characters long. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Username) > 32 {
		errorBuffer.WriteString("Username should be at most 32 characters long. \n")
	}

	if utf8.RuneCountInString(unregisteredUser.Password) == 0 {
		errorBuffer.WriteString("Password is a required field. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Password) < 6 {
		errorBuffer.WriteString("Password should be at least 6 characters long. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Password) > 64 {
		errorBuffer.WriteString("Password should be at most 64 characters long. \n")
	}

	strongPassword, _ := regexp.MatchString(".*[a-z].*([A-Z].*[0-9]|[0-9].*[A-Z])|.*[A-Z]" +
		".*([a-z].*[0-9]|[0-9].*[a-z])|.*[0-9].*([a-z].*[A-Z]|[A-Z].*[a-z])", unregisteredUser.Password)

	if strongPassword == false {
		errorBuffer.WriteString("Password should contain at least one number, " +
			"one capital character and one lowercase character. \n")
	}

	if utf8.RuneCountInString(unregisteredUser.Email) == 0 {
		errorBuffer.WriteString("Email is a required field. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Email) < 6 {
		errorBuffer.WriteString("Email should be at least 6 characters long. \n")
	} else if utf8.RuneCountInString(unregisteredUser.Email) > 128 {
		errorBuffer.WriteString("Email should be at most 128 characters long. \n")
	}

	validEmail, _ := regexp.MatchString(`[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}`, unregisteredUser.Email)

	if validEmail == false {
		errorBuffer.WriteString("Email is not in a valid format. \n")
	}

	if len(errorBuffer.String()) > 0 {
		return errors.New(errorBuffer.String())
	} else {
		return nil
	}
}
