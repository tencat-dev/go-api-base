package infra

import (
	"github.com/goforj/wire"

	"github.com/tencat-dev/go-api-base/internal/infra/auth"
)

// ProviderSetInfra is infra providers.
var ProviderSetInfra = wire.NewSet(
	auth.NewJWTMaker,
	auth.NewTokenMaker,
	auth.TokenParser,
	auth.NewPasswordHasher,
)
