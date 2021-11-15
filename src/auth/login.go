package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	db "github.com/options-go/src/database"
	"github.com/options-go/src/users"
	"github.com/options-go/src/utils"
)

type loginData struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type DBToken struct {
	ID           uuid.UUID `gorm:"primaryKey" json:"id"`
	Email        string    `json:"email"`
	Valid        bool      `json:"valid"`
	RefreshToken uuid.UUID `json:"refresh_token"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:nano"`
}

func AuthenticateUser(res http.ResponseWriter, req *http.Request) {
	loginCreds := loginData{}
	body, read_err := ioutil.ReadAll(req.Body)
	unmarshal_error := json.Unmarshal(body, &loginCreds)

	done_token := make(chan bool)

	token_id := uuid.New()
	refresh_token := uuid.New()

	if read_err != nil && unmarshal_error != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("bad request"))
		return
	}
	user := users.User{
		Email: loginCreds.Email,
		Pass:  utils.SHA256(loginCreds.Pass),
	}

	go AddToken(done_token, user.Email, token_id, refresh_token)

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

	token, token_error := GenerateToken(dbData.FirstName, dbData.LastName, dbData.Email, dbData.ID, token_id, refresh_token)
	if token_error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(token_error.Error()))
	}

	token_success := <-done_token

	if !token_success {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("could not generate token"))
		return
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
		Find(&stringInterface, "email = ? AND pass = ?", user.Email, user.Pass).
		Table("tokens")
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

func AddToken(done chan bool, email string, tid uuid.UUID, rft uuid.UUID) {
	db_token := DBToken{
		ID:           tid,
		Email:        email,
		Valid:        true,
		RefreshToken: rft,
	}

	results := db.GormCon().
		Table("tokens").
		Create(&db_token)

	if results.Error != nil {
		log.Println(results.Error.Error())
		done <- false
		return
	}

	if results.RowsAffected == 0 {
		log.Println("no rows affected")
		done <- false
		return
	}

	done <- true
	return
}
