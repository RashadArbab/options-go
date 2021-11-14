package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Payload struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	User_ID      string    `json:"user_id"`
	Token_ID     string    `json:"token_id"`
	RefreshToken string    `json:"refresh_token"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type JWTKey struct {
	secretKey string
}

// there are two options for security, create a very short token 5 mintues
// or check if token id is revoked.

// i suppose one thing to do would be that if a refresh token is used mutltiple times
// block all further refresh tokens associated with a user. That will force them to relogin

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("token expired")
	}
	return nil
}

type JWT interface {
	GenerateToken(fn string, ln string, em string, uid string) (string, error)
	VerifyToken(token string)
	generatePayload(fn string, ln string, em string, uid string) (*Payload, error)
}

func generatePayload(fn string, ln string, em string, uid string) (*Payload, error) {
	payload := Payload{
		FirstName:    fn,
		LastName:     ln,
		Email:        em,
		User_ID:      uid,
		Token_ID:     uuid.NewString(),
		RefreshToken: uuid.NewString(),
		IssuedAt:     time.Now(),
		ExpiredAt:    time.Now().Add(time.Hour),
	}

	return &payload, nil
}

func GenerateToken(fn string, ln string, em string, uid string) (string, error) {
	payload, payload_err := generatePayload(fn, ln, em, uid)
	if payload_err != nil {
		return "", payload_err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signingKey := os.Getenv("TOKEN_KEY")
	log.Println(signingKey)
	signedToken, sign_err := jwtToken.SignedString([]byte(signingKey))
	if sign_err != nil {
		return "", sign_err
	}
	return signedToken, nil
}
