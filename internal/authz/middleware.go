package authz

import (
	"context"
	"strings"

	"github.com/anhnmt/go-authxx/token"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"

	"github.com/tencat-dev/go-api-base/internal/biz"
)

type AuthzMiddleware middleware.Middleware

func NewAuthzMiddleware(
	pc biz.PermissionChecker,
	r *AuthzRegistry,
	userRepo biz.UserRepo,
	parser token.TokenParser,
) AuthzMiddleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			header, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("NO_TRANSPORT", "no transport")
			}

			fullMethod := header.Operation()

			// 🔥 Fast path: public API
			perm, exists := r.Get(fullMethod)
			if !exists {
				return next(ctx, req)
			}

			// 🔥 Protected API → verify JWT first
			tokenStr, err := extractBearerToken(header)
			if err != nil {
				return nil, err
			}

			claims, err := parseTokenOfType(parser, tokenStr, token.AccessToken)
			if err != nil {
				return nil, err
			}

			userID, err := validateUserFromToken(ctx, userRepo, claims)
			if err != nil {
				return nil, err
			}

			if err = checkPermission(userID, pc, perm.Object, perm.Action); err != nil {
				return nil, err
			}

			ctx = token.NewContext(ctx, claims)
			return next(ctx, req)
		}
	}
}

func extractBearerToken(header transport.Transporter) (string, error) {
	tokenKey := header.RequestHeader().Get(authorizationKey)
	parts := strings.SplitN(tokenKey, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], bearerWord) {
		return "", ErrMissingJwtToken
	}

	return parts[1], nil
}

func parseTokenOfType(parser token.TokenParser, tokenStr string, expectedType token.TokenType) (*token.TokenPayload, error) {
	payload, err := parser.ParseToken(tokenStr)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrExpiredToken):
			return nil, ErrTokenExpired
		default:
			return nil, ErrTokenInvalid
		}
	}

	if payload.Type != expectedType {
		return nil, ErrTokenInvalid
	}

	return payload, nil
}

func validateUserFromToken(
	ctx context.Context,
	userRepo biz.UserRepo,
	claims *token.TokenPayload,
) (uuid.UUID, error) {
	user, err := userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return uuid.Nil, err
	}

	if claims.TokenVersion != user.TokenVersion {
		return uuid.Nil, errors.Unauthorized("INVALID_TOKEN", "token revoked")
	}

	return claims.UserID, nil
}

func checkPermission(
	userID uuid.UUID,
	pc biz.PermissionChecker,
	object, action string,
) error {
	allowed, err := pc.Can(userID, object, action)
	if err != nil {
		return err
	}

	if !allowed {
		return errors.Forbidden("ACCESS_DENIED", "permission denied")
	}

	return nil
}
