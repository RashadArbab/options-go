package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func SHA256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))

}

func SHA512(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}
