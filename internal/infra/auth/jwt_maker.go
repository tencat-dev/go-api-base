package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/conf"
)

type JWTClaims struct {
	SessionID    string        `json:"sid,omitempty"`
	Type         biz.TokenType `json:"type,omitempty"`
	TokenVersion int32         `json:"token_version,omitempty"`
	jwt.RegisteredClaims
}

type jwtMaker struct {
	secretKey string
}

func NewJWTMaker(c *conf.JWT) biz.TokenMaker {
	return &jwtMaker{secretKey: c.Secret}
}

func (j *jwtMaker) CreateToken(
	now time.Time,
	p biz.TokenPayload,
) (string, error) {
	claims := JWTClaims{
		SessionID:    p.SessionID.String(),
		Type:         p.Type,
		TokenVersion: p.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        p.TokenID.String(),
			Subject:   p.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
