package data

import (
	"github.com/anhnmt/go-authxx/rbac"
	"github.com/casbin/casbin/v3"
)

func NewCasbinAuthz(enforcer casbin.IEnforcer) *rbac.CasbinAuthz {
	return rbac.New(enforcer)
}

func NewPermissionChecker(c *rbac.CasbinAuthz) rbac.Checker {
	return c
}

func NewPermissionManager(c *rbac.CasbinAuthz) rbac.Manager {
	return c
}

func NewCasbinEnforcer(data *Data) (casbin.IEnforcer, error) {
	return rbac.NewCasbinEnforcer(data.db,
		rbac.WithModelPath("configs/rbac_model.conf"),
		rbac.WithIndex("ptype", "v0", "v1", "v2"), // policy: sub, obj, act
		rbac.WithIndex("ptype", "v0", "v1"),       // grouping: user -> role
	)
}
