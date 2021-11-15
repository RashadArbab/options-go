package main

import (
	"fmt"

	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/options-go/src/auth"
	"github.com/options-go/src/blackscholes"
	db "github.com/options-go/src/database"
	users "github.com/options-go/src/users"
)

func main() {
	fmt.Println("starting server on port:8080")

	env_err := godotenv.Load(".env")
	if env_err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		fmt.Println("Loaded env file")
	}

	_, mongo_err := db.MongoConnect()
	if mongo_err != nil {
		log.Fatal("could not connect to db")
	}

	db.GormCon()

	http.HandleFunc("/qtradeOauth", handleQtrade)
	http.HandleFunc("/calculate", blackscholes.CalcBlackScholes)
	http.HandleFunc("/users", users.UserEndpoints)
	http.HandleFunc("/positions", auth.AuthRequired(users.PostionsEndpoints))
	http.HandleFunc("/position", auth.AuthRequired(users.PositionEndpoints))
	http.HandleFunc("/login", auth.AuthenticateUser)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleRed(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("we have been redirected"))
}

func handleQtrade(res http.ResponseWriter, req *http.Request) {
	client_id := "V9UGAlV9xrCXTAAO09KAzP2EehobEg"

	url := "https://login.questrade.com/oauth2/authorize?client_id=" + client_id + "&response_type=code&redirect_uri=https://fufvciexl3.execute-api.us-east-1.amazonaws.com/default/redirectLambda"
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}
