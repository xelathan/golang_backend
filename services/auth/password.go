package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckHashedPassword(inputPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), ([]byte(inputPassword)))
}
