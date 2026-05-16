package bcrypt

import "golang.org/x/crypto/bcrypt"

// ToHash -.
func ToHash(value string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
}

// CompareHash -.
func CompareHash(v1, v2 []byte) error {
	return bcrypt.CompareHashAndPassword(v1, v2)
}
