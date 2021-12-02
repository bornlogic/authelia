package authentication

import "golang.org/x/crypto/bcrypt"

// CheckPassword check a password against a hash.
func CheckPassword(hashed, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return false, err
	}
	return true, nil
}

// HashPassword generate a hashed password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}
