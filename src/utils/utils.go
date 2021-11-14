package utils

import (
	"errors"
	"net/http"
)

func VerifyHTTPMethod(method string, req *http.Request, res http.ResponseWriter) bool {
	if req.Method != method {
		res.WriteHeader(404)
		res.Write([]byte("Not Found!"))
		return false
	} else {
		return true
	}

}

func VerifyParams(fields []string, res http.ResponseWriter, req *http.Request) error {
	for _, element := range fields {
		if !req.URL.Query().Has(element) {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("missing parameters"))
			return errors.New("missing parameter")
		}
	}
	return nil
}
