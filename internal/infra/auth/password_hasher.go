package auth

import (
	"github.com/anhnmt/go-pkgxx/password"
)

func NewPasswordHasher() password.Hasher {
	return password.NewArgon2Hasher()
}
