package auth

import (
	"log"
	"net/http"
	"strings"
)

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		bearerToken := req.Header.Get("Authorization")
		token := strings.Split(bearerToken, "Bearer")
		trim_token := strings.Trim(token[1], " ")
		payload, err := VerifyToken(trim_token)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte(err.Error()))
			return
		}

		log.Println(payload)

		next.ServeHTTP(res, req)

	}

}
