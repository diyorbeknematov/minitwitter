package security

import "golang.org/x/crypto/bcrypt"

func HashToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)

	return string(hashedToken), err
}

func VerifyToken(token, hashToken string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashToken), []byte(token)) == nil
}
