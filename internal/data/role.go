package data

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
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
	setters := lo.Map(roles, func(role *biz.Role, index int) *models.RoleSetter {
		return &models.RoleSetter{
			Name:        omit.From(role.Name),
			Description: omitnull.From(role.Description),
			IsSystem:    omitnull.From(role.IsSystem),
		}
	})

	_, err := models.Roles.Insert(
		bob.ToMods(setters...),
		im.OnConflict("name").DoNothing(),
	).Exec(ctx, r.data.db)
	if err != nil {
		return err
	}

	return nil
}

func (r *roleRepo) ExistByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exist, err := models.Roles.Query(
		sm.Where(models.Roles.Columns.ID.EQ(psql.Arg(id))),
	).Exists(ctx, r.data.db)
	if err != nil {
		return exist, err
	}

	return exist, nil
}
