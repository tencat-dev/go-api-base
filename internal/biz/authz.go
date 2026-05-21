package biz

import (
	"context"
	"fmt"

	"github.com/anhnmt/go-authxx/rbac"
	"github.com/google/uuid"
)

// AuthzBiz is a Auth usecase.
type AuthzBiz struct {
	pm       rbac.Manager
	userRepo UserRepo
}

// NewAuthzBiz new a Auth usecase.
func NewAuthzBiz(pm rbac.Manager, userRepo UserRepo) *AuthzBiz {
	return &AuthzBiz{
		pm:       pm,
		userRepo: userRepo,
	}
}

func (b *AuthzBiz) GrantRole(ctx context.Context, userID uuid.UUID, role string) error {
	userExist, err := b.userRepo.ExistByID(ctx, userID)
	if err != nil {
		return err
	}

	if !userExist {
		return fmt.Errorf("user not exist")
	}

	return b.pm.GrantRole(userID.String(), role)
}

func (b *AuthzBiz) RevokeRole(ctx context.Context, userID uuid.UUID, role string) error {
	userExist, err := b.userRepo.ExistByID(ctx, userID)
	if err != nil {
		return err
	}

	if !userExist {
		return fmt.Errorf("user not exist")
	}

	return b.pm.RevokeRole(userID.String(), role)
}

func (b *AuthzBiz) GrantPermission(
	_ context.Context,
	role, object, action string,
) error {
	return b.pm.GrantPermission(role, object, action)
}
