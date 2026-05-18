package auth

import (
	"github.com/anhnmt/go-authxx/password"
)

func NewPasswordHasher() password.Hasher {
	return password.NewArgon2Hasher()
}
