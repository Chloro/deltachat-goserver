package endpoints

import(
	"deltachat/models"
	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func GetAllUsers(request *http.Request, response http.ResponseWriter, s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	collection := session.DB("deltachat").C("users")

	var users []models.User

	err := collection.Find(bson.M{}).All(&users)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(200)
	json.NewEncoder(response).Encode(users)
}
