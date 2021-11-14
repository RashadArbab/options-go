package users

import (
	"encoding/json"
	"net/http"
	"strings"

	db "github.com/options-go/src/database"
	"github.com/options-go/src/utils"
)

func GetUser(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	userID := params.Get("user_id")
	param_err := utils.VerifyParams([]string{"user_id"}, res, req)
	if param_err != nil {
		return
	}

	stringInterface := map[string]interface{}{}

	db.GormCon().
		Table("users").
		Select("users.id", "users.first_name", "users.last_name", "users.email").
		Where(&User{ID: userID}).
		Scan(&stringInterface)

	dataJson, err := json.MarshalIndent(stringInterface, "", "    ")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write(dataJson)

}

func GetPositions(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	userID := params.Get("user_id")
	param_err := utils.VerifyParams([]string{"user_id"}, res, req)
	if param_err != nil {
		return
	}
	stringInterface := map[string]interface{}{}

	results := db.GormCon().
		Table("positions").
		Find(&stringInterface, "user_id = ?", userID)
	if results.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(results.Error.Error()))
		return
	}

	dataJson, err := json.MarshalIndent(stringInterface, "", "    ")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	if results.RowsAffected == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(dataJson)
}

func GetPosition(res http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	userID := params.Get("user_id")
	ticker := params.Get("ticker")
	param_err := utils.VerifyParams([]string{"user_id", "ticker"}, res, req)
	if param_err != nil {
		return
	}

	ticker = strings.ToUpper(ticker)
	stringInterface := map[string]interface{}{}

	results := db.GormCon().
		Table("positions").
		Find(&stringInterface, "user_id = ? AND ticker = ?", userID, ticker)
	if results.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(results.Error.Error()))
		return
	}

	if results.RowsAffected == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	dataJson, err := json.MarshalIndent(stringInterface, "", "    ")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(dataJson)
}
