package biz

import (
	"context"

	"github.com/google/uuid"
)

// Permission is a Permission model.
type Permission struct {
	ID     uuid.UUID `json:"id,omitempty"`
	Object string    `json:"object,omitempty"`
	Action string    `json:"action,omitempty"`
}

// PermissionRepo is a Greater repo.
type PermissionRepo interface {
	Saves(context.Context, []*Permission) error
	ExistByID(context.Context, uuid.UUID) (bool, error)
}

// PermissionBiz is a Permission usecase.
type PermissionBiz struct {
	repo PermissionRepo
}

// NewPermissionBiz new a Permission usecase.
func NewPermissionBiz(repo PermissionRepo) *PermissionBiz {
	return &PermissionBiz{repo: repo}
}

// CreatePermissions creates a Permission, and returns the new Permission.
func (b *PermissionBiz) CreatePermissions(ctx context.Context, pers []*Permission) error {
	return b.repo.Saves(ctx, pers)
}
