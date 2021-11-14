package users

import (
	"net/http"
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Pass      string    `pg:"pass"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:nano"`
}

type Position struct {
	ID     string `gorm:"primaryKey" json:"id"`
	UserID string `json:"user_id" gorm:"foreignKey:UserID,index:ticker_idx"`
	Ticker string `gorm:"index:ticker_idx" json:"ticker"`
	Amount int    `pg:"amount" json:"amount"`
}

func UserEndpoints(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		GetUser(res, req)
	} else if req.Method == "POST" {
		CreateUser(res, req)
	} else {
		res.WriteHeader(404)
		res.Write([]byte("Not Found!"))
		return
	}
}

func PositionEndpoints(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		GetPosition(res, req)
	} else if req.Method == "POST" {
		CreatePositon(res, req)
	} else if req.Method == "DELETE" {
		DeletePositon(res, req)
	} else if req.Method == "PUT" {
		UpdatePositon(res, req)
	} else {
		res.WriteHeader(404)
		res.Write([]byte("Not Found!"))
		return
	}
}

func PostionsEndpoints(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		GetPositions(res, req)
	} else {
		res.WriteHeader(404)
		res.Write([]byte("Not Found!"))
		return
	}
}
