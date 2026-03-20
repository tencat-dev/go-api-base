package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Role is a Role model.
type Role struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// RoleRepo is a Greater repo.
type RoleRepo interface {
	Saves(context.Context, []*Role) error
	ExistByID(context.Context, uuid.UUID) (bool, error)
}

// RoleBiz is a Role usecase.
type RoleBiz struct {
	repo RoleRepo
}

// NewRoleBiz new a Role usecase.
func NewRoleBiz(repo RoleRepo) *RoleBiz {
	return &RoleBiz{repo: repo}
}

// CreateRoles creates a Role, and returns the new Role.
func (b *RoleBiz) CreateRoles(ctx context.Context, r []*Role) error {
	return b.repo.Saves(ctx, r)
}
