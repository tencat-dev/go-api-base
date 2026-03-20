package authz

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/conf"
	"github.com/tencat-dev/go-api-base/internal/infra/auth"
)

type AuthzMiddleware middleware.Middleware

func NewAuthzMiddleware(
	jwtConf *conf.JWT,
	pc biz.PermissionChecker,
	r *AuthzRegistry,
	userRepo biz.UserRepo,
) AuthzMiddleware {
	keyFunc := func(*jwtv5.Token) (any, error) {
		return []byte(jwtConf.Secret), nil
	}

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

			claims, err := parseAndValidateJWT(tokenStr, keyFunc)
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

			ctx = NewContext(ctx, claims)
			return next(ctx, req)
		}
	}
}

func extractBearerToken(header transport.Transporter) (string, error) {
	token := header.RequestHeader().Get(authorizationKey)
	parts := strings.SplitN(token, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], bearerWord) {
		return "", ErrMissingJwtToken
	}

	return parts[1], nil
}

func parseAndValidateJWT(tokenStr string, keyFunc jwtv5.Keyfunc) (*auth.JWTClaims, error) {
	tokenInfo, err := jwtv5.ParseWithClaims(tokenStr, &auth.JWTClaims{}, keyFunc)
	if err != nil {
		switch {
		case errors.Is(err, jwtv5.ErrTokenMalformed),
			errors.Is(err, jwtv5.ErrTokenUnverifiable):
			return nil, ErrTokenInvalid
		case errors.Is(err, jwtv5.ErrTokenExpired),
			errors.Is(err, jwtv5.ErrTokenNotValidYet):
			return nil, ErrTokenExpired
		default:
			return nil, ErrTokenParseFail
		}
	}

	if !tokenInfo.Valid {
		return nil, ErrTokenInvalid
	}

	if tokenInfo.Method != jwtv5.SigningMethodHS256 {
		return nil, ErrUnSupportSigningMethod
	}

	claims, ok := tokenInfo.Claims.(*auth.JWTClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	if claims.Type != biz.AccessToken {
		return nil, errors.Unauthorized("NO_ACCESS_TOKEN", "access token required")
	}

	return claims, nil
}

func validateUserFromToken(
	ctx context.Context,
	userRepo biz.UserRepo,
	claims *auth.JWTClaims,
) (uuid.UUID, error) {
	sub, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, errors.Unauthorized("INVALID_TOKEN", err.Error())
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, errors.Unauthorized("INVALID_TOKEN", err.Error())
	}

	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}

	if claims.TokenVersion != user.TokenVersion {
		return uuid.Nil, errors.Unauthorized("INVALID_TOKEN", "token revoked")
	}

	return userID, nil
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
