package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	db "github.com/options-go/src/database"
)

type Payload struct {
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	User_ID      string    `json:"user_id"`
	Token_ID     uuid.UUID `json:"token_id"`
	RefreshToken uuid.UUID `json:"refresh_token"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type JWTKey struct {
	secretKey string
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

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

func TokenRevoked(payload Payload) error {

	token := DBToken{}

	results := db.GormCon().
		Table("tokens").
		Find(&token, "id = ?", payload.Token_ID)

	if results.Error != nil {
		return results.Error
	}

	if !token.Valid {
		return errors.New("token not valid")
	}
	return nil

}

type JWT interface {
	GenerateToken(fn string, ln string, em string, uid string) (string, error)
	VerifyToken(token string)
	generatePayload(fn string, ln string, em string, uid string) (*Payload, error)
}

func generatePayload(fn string, ln string, em string, uid string, tid uuid.UUID, rft uuid.UUID) (*Payload, error) {
	payload := Payload{
		FirstName:    fn,
		LastName:     ln,
		Email:        em,
		User_ID:      uid,
		Token_ID:     tid,
		RefreshToken: rft,
		IssuedAt:     time.Now(),
		ExpiredAt:    time.Now().Add(time.Hour),
	}

	return &payload, nil
}

func GenerateToken(fn string, ln string, em string, uid string, tid uuid.UUID, rft uuid.UUID) (string, error) {
	payload, payload_err := generatePayload(fn, ln, em, uid, tid, rft)
	if payload_err != nil {
		return "", payload_err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signingKey := os.Getenv("TOKEN_KEY")
	signedToken, sign_err := jwtToken.SignedString([]byte(signingKey))
	if sign_err != nil {
		return "", sign_err
	}
	return signedToken, nil
}

//This is the overall TokenVerify Method
func VerifyToken(token string) (*Payload, error) {
	//first we need to create the decode keyFunc and we are checking that
	// the method of signing is the same as what we used
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("TOKEN_KEY")), nil
	}

	// we use the keyfunc to to check if token is valid interms of time
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// parse payload to payload object
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	//validate that token is not rejected in our database
	revoked_err := TokenRevoked(*payload)
	if revoked_err != nil {
		//implement remove all tokens from user.

		return nil, ErrInvalidToken
	}

	return payload, nil
}
