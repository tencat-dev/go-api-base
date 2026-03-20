package biz

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PermissionChecker interface {
	Can(userId uuid.UUID, obj, act string) (bool, error)
}

type PermissionManager interface {
	GrantRole(userID uuid.UUID, role string) error
	RevokeRole(userID uuid.UUID, role string) error
	GrantPermission(role, object, action string) error
	GrantPermissions(rules [][]string) error
	RevokePermission(role, object, action string) error
	DeleteRolesForUser(userID uuid.UUID) error
}

// AuthzBiz is a Auth usecase.
type AuthzBiz struct {
	pm       PermissionManager
	userRepo UserRepo
}

// NewAuthzBiz new a Auth usecase.
func NewAuthzBiz(pm PermissionManager, userRepo UserRepo) *AuthzBiz {
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

	return b.pm.GrantRole(userID, role)
}

func (b *AuthzBiz) RevokeRole(ctx context.Context, userID uuid.UUID, role string) error {
	userExist, err := b.userRepo.ExistByID(ctx, userID)
	if err != nil {
		return err
	}

	if !userExist {
		return fmt.Errorf("user not exist")
	}

	return b.pm.RevokeRole(userID, role)
}

func (b *AuthzBiz) GrantPermission(
	_ context.Context,
	role, object, action string,
) error {
	return b.pm.GrantPermission(role, object, action)
}
