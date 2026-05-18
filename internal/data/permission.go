package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/infra/persistence/postgres"
)

type permissionRepo struct {
	data *Data
	log  *log.Helper
}

// NewPermissionRepo .
func NewPermissionRepo(data *Data, logger *log.Helper) biz.PermissionRepo {
	return &permissionRepo{
		data: data,
		log:  logger,
	}
}

func (r *permissionRepo) Saves(ctx context.Context, pers []*biz.Permission) error {
	tx, err := r.data.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.data.queries.WithTx(tx)

	for _, per := range pers {
		err := qtx.InsertPermission(ctx, postgres.InsertPermissionParams{
			Object: per.Object,
			Action: per.Action,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *permissionRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := r.data.queries.ExistsPermission(ctx, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
