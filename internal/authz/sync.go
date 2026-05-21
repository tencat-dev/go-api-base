package authz

import (
	"context"

	"github.com/anhnmt/go-authxx/rbac"

	authzv1 "github.com/tencat-dev/go-api-base/api/authz/v1"
	"github.com/tencat-dev/go-api-base/internal/biz"
)

type SeederFromRegistry struct {
	pm            rbac.Manager
	registry      *AuthzRegistry
	roleBiz       *biz.RoleBiz
	permissionBiz *biz.PermissionBiz
}

func NewSeederFromRegistry(pm rbac.Manager, r *AuthzRegistry, roleBiz *biz.RoleBiz,
	permissionBiz *biz.PermissionBiz) *SeederFromRegistry {
	return &SeederFromRegistry{
		pm:            pm,
		registry:      r,
		roleBiz:       roleBiz,
		permissionBiz: permissionBiz,
	}
}

func (s *SeederFromRegistry) AutoSync() error {
	m := s.registry.data.Load().(map[string]*authzv1.PermissionOption)

	// Ước lượng capacity để giảm re-alloc
	policies := make([][]string, 0, len(m)*2)

	// dedup trong memory (tránh duplicate trong cùng 1 proto load)
	seen := make(map[string]struct{})
	roles := make([]*biz.Role, 0, 2)
	permissions := make([]*biz.Permission, 0, len(m))

	for _, perm := range m {
		permissions = append(permissions, &biz.Permission{
			Object: perm.Object,
			Action: perm.Action,
		})

		for _, role := range perm.Roles {
			roles = append(roles, &biz.Role{
				Name:        role,
				IsSystem:    true,
				Description: "Default role",
			})

			key := role + "|" + perm.Object + "|" + perm.Action
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			policies = append(policies, []string{
				role,
				perm.Object,
				perm.Action,
			})
		}
	}

	ctx := context.Background()

	if len(roles) > 0 {
		err := s.roleBiz.CreateRoles(ctx, roles)
		if err != nil {
			return err
		}
	}

	if len(permissions) > 0 {
		err := s.permissionBiz.CreatePermissions(ctx, permissions)
		if err != nil {
			return err
		}
	}

	if len(policies) > 0 {
		// 🔥 Batch + ignore duplicate DB entries
		err := s.pm.GrantPermissions(policies)
		if err != nil {
			return err
		}
	}

	return nil
}
