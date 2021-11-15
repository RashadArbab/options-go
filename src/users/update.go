package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	db "github.com/options-go/src/database"
	"github.com/options-go/src/utils"
	"gorm.io/gorm"
)

func UpdatePositon(res http.ResponseWriter, req *http.Request) {

	log.Println("create-position:")

	body, read_err := utils.ReadRequest(res, req) //read into []byte
	if read_err != nil {
		return
	}
	fmt.Println(body)
	var pos Position
	unmarshal_err := json.Unmarshal(body, &pos) //marshal into struct Position
	if utils.HandleStandardError(unmarshal_err, res) != nil {
		return
	}
	pos.Ticker = strings.ToUpper(pos.Ticker)

	result := db.GormCon().
		Table("positions").
		Where("user_id = ?", pos.UserID).
		Where("ticker = ?", pos.Ticker).
		UpdateColumn("amount", gorm.Expr("amount + ?", pos.Amount))

	if result.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(result.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Successfully updated position"))

}
