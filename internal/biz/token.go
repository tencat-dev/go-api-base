package biz

import (
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(now time.Time, payload TokenPayload) (string, error)
}

type TokenPayload struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	TokenID      uuid.UUID
	Type         TokenType
	TokenVersion int32
	TTL          time.Duration
}

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)
