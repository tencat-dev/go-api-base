package data

import (
	"context"

	"github.com/google/uuid"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/infra/persistence/postgres"

	"github.com/go-kratos/kratos/v2/log"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo .
func NewRoleRepo(data *Data, logger *log.Helper) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  logger,
	}
}

func (r *roleRepo) Saves(ctx context.Context, roles []*biz.Role) error {
	tx, err := r.data.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.data.queries.WithTx(tx)

	for _, role := range roles {
		err := qtx.InsertRole(ctx, postgres.InsertRoleParams{
			Name:        role.Name,
			Description: new(role.Description),
			IsSystem:    new(role.IsSystem),
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *roleRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := r.data.queries.ExistsRole(ctx, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
