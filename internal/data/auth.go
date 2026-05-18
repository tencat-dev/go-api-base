package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/tencat-dev/go-api-base/internal/biz"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuthRepo .
func NewAuthRepo(data *Data, logger *log.Helper) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  logger,
	}
}

func (r *authRepo) FindByEmail(ctx context.Context, email string) (*biz.Auth, error) {
	auth, err := r.data.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &biz.Auth{
		ID:           auth.ID,
		Name:         auth.Name,
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash,
	}, nil
}
