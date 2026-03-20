package data

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"

	"github.com/tencat-dev/go-api-base/internal/biz"
	"github.com/tencat-dev/go-api-base/internal/infra/persistence/postgres/models"

	"github.com/go-kratos/kratos/v2/log"
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
	setters := lo.Map(pers, func(per *biz.Permission, _ int) *models.PermissionSetter {
		return &models.PermissionSetter{
			Object: omit.From(per.Object),
			Action: omit.From(per.Action),
		}
	})

	_, err := models.Permissions.Insert(
		bob.ToMods(setters...),
		im.OnConflict("object", "action").DoNothing(),
	).Exec(ctx, r.data.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *permissionRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exist, err := models.Permissions.Query(
		sm.Where(models.Permissions.Columns.ID.EQ(psql.Arg(id))),
	).Exists(ctx, r.data.db)
	if err != nil {
		return exist, err
	}

	return exist, nil
}
