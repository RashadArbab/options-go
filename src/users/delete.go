package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	db "github.com/options-go/src/database"
	"github.com/options-go/src/utils"
)

func DeletePositon(res http.ResponseWriter, req *http.Request) {

	body, read_err := utils.ReadRequest(res, req) //read into []byte
	if read_err != nil {
		return
	}
	var pos Position
	unmarshal_err := json.Unmarshal(body, &pos) //marshal into struct Position
	if utils.HandleStandardError(unmarshal_err, res) != nil {
		return
	}
	pos.Ticker = strings.ToUpper(pos.Ticker)
	pos.ID = uuid.NewString()

	result := db.GormCon().
		Table("positions").
		Where("user_id = ? AND ticker = ?", pos.UserID, pos.Ticker).
		Delete(&Position{})
	if result.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(result.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Successfully deleted position"))

}
