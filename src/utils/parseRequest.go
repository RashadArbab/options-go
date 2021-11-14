package utils

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ReadRequest(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error reading body")
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(err.Error()))
		return nil, err
	}
	return body, nil

}
