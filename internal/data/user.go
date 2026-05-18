package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/infra/persistence/postgres"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger *log.Helper) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  logger,
	}
}

func (r *userRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	tokenVersion := int32(0)
	if u.TokenVersion > 0 {
		tokenVersion = u.TokenVersion
	}

	insertedUser, err := r.data.queries.InsertUser(ctx, postgres.InsertUserParams{
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		TokenVersion: tokenVersion,
	})
	if err != nil {
		return nil, err
	}

	u.ID = insertedUser.ID
	u.CreatedAt = insertedUser.CreatedAt
	u.UpdatedAt = insertedUser.UpdatedAt

	return u, nil
}

func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	updatedUser, err := r.data.queries.UpdateUser(ctx, postgres.UpdateUserParams{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
	if err != nil {
		return nil, err
	}

	u.UpdatedAt = updatedUser.UpdatedAt
	return u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.User, error) {
	user, err := r.data.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &biz.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		TokenVersion: user.TokenVersion,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *userRepo) ListAll(ctx context.Context) ([]*biz.User, error) {
	userSlice, err := r.data.queries.ListUsers(ctx, int32(10))
	if err != nil {
		return nil, err
	}

	var users []*biz.User
	for _, user := range userSlice {
		users = append(users, &biz.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return users, nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := r.data.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := r.data.queries.ExistsUser(ctx, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
