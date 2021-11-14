package utils

import "net/http"

func HandleStandardError(e error, res http.ResponseWriter) error {
	if e != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(e.Error()))
		return e
	}
	return nil
}
