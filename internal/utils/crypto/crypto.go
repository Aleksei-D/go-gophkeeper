package crypto

import (
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/credentials"
)

// Hasher интерфейс работы с хешем
//
//go:generate mockgen -destination=mocks/mock_hasher.go -package=mocks . Hasher
type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// BcryptHasher работа с хешем
type BcryptHasher struct{}

// HashPassword возврат хеша от пароля
func (b *BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверка пароля
func (b *BcryptHasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateTLSCreds(certFile string) (credentials.TransportCredentials, error) {
	return credentials.NewClientTLSFromFile(certFile, "")
}
