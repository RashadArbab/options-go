package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	db "github.com/options-go/src/database"
	"github.com/options-go/src/users"
	"github.com/options-go/src/utils"
)

type loginData struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func AuthenticateUser(res http.ResponseWriter, req *http.Request) {
	loginCreds := loginData{}
	body, read_err := ioutil.ReadAll(req.Body)
	unmarshal_error := json.Unmarshal(body, &loginCreds)
	if read_err != nil && unmarshal_error != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("bad request"))
		return
	}
	user := users.User{
		Email: loginCreds.Email,
		Pass:  utils.SHA512(loginCreds.Pass),
	}
	user_data, statusCode, err := FindUser(user)
	if err != nil {
		res.WriteHeader(statusCode)
		res.Write([]byte(err.Error()))
		return
	}
	var dbData = users.User{}
	unmarshal_err := json.Unmarshal(user_data, &dbData)
	if unmarshal_err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(unmarshal_err.Error()))
	}

	token, token_error := GenerateToken(dbData.FirstName, dbData.LastName, dbData.Email, dbData.ID)
	if token_error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(token_error.Error()))
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(token))
	return
}

func FindUser(user users.User) ([]byte, int, error) {
	stringInterface := map[string]interface{}{}
	results := db.GormCon().
		Table("users").
		Select("id", "first_name", "last_name", "email").
		Find(&stringInterface, "email = ? AND pass = ?", user.Email, user.Pass)

	if results.Error != nil {
		return nil, http.StatusInternalServerError, results.Error
	}

	if len(stringInterface) == 0 {
		return nil, http.StatusUnauthorized, errors.New("unauthorized")
	}

	dataJson, err := json.MarshalIndent(stringInterface, "", "    ")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return dataJson, http.StatusOK, nil

}
