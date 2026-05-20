package auth

import (
	"github.com/anhnmt/go-authxx/token"

	"github.com/tencat-dev/go-api-base/internal/conf"
)

func NewJWTMaker(c *conf.JWT) (*token.JWTMaker, error) {
	return token.NewJWTMaker(c.Secret)
}

func NewTokenMaker(jwt *token.JWTMaker) token.TokenMaker {
	return jwt
}

func TokenParser(jwt *token.JWTMaker) token.TokenParser {
	return jwt
}
