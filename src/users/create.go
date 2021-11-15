package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	db "github.com/options-go/src/database"
	"github.com/options-go/src/utils"
	"gorm.io/gorm/clause"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {

	log.Println("create-user:")
	body, read_err := ioutil.ReadAll(req.Body)
	if read_err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Request format error."))
		return
	}

	var user = User{}
	unmarshal_err := json.Unmarshal(body, &user)
	if unmarshal_err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Internal Server Error"))
		return
	}

	user.Pass = utils.SHA256(user.Pass)
	user.ID = uuid.NewString()

	result := db.GormCon().Create(&user)
	if result.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(result.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Successfully created user"))

}

func CreatePositon(res http.ResponseWriter, req *http.Request) {

	log.Println("create-position:")
	body, read_err := utils.ReadRequest(res, req) //read into []byte
	if read_err != nil {
		return
	}
	fmt.Println(string(body))
	var pos Position
	unmarshal_err := json.Unmarshal(body, &pos) //marshal into struct Position
	if utils.HandleStandardError(unmarshal_err, res) != nil {
		return
	}
	pos.Ticker = strings.ToUpper(pos.Ticker)
	pos.ID = uuid.NewString()

	result := db.GormCon().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "ticker"}},
		//change to throw error position already exists
		DoUpdates: clause.AssignmentColumns([]string{"amount"}),
	}).Create(&pos)

	if result.Error != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(result.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("Successfully created position"))

}
