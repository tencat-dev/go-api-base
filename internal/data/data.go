package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/goforj/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tencat-dev/go-api-base/internal/conf"
	"github.com/tencat-dev/go-api-base/internal/infra/persistence/postgres"
)

// ProviderSetData is data providers.
var ProviderSetData = wire.NewSet(
	NewData,
	NewCasbinEnforcer,
	NewCasbinAuthz,
	NewPermissionChecker,
	NewPermissionManager,
	NewUserRepo,
	NewAuthRepo,
	NewRoleRepo,
	NewPermissionRepo,
)

// Data wraps database client.
type Data struct {
	db      *pgxpool.Pool
	queries *postgres.Queries
}

// NewData creates a new Data instance with PostgreSQL connection.
func NewData(ctx context.Context, c *conf.DatabaseConfig, logHelper *log.Helper) (*Data, func(), error) {
	if c == nil {
		return nil, nil, fmt.Errorf("database configuration is missing")
	}

	config, err := pgxpool.ParseConfig(c.Dsn)
	if err != nil {
		return nil, nil, err
	}
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.PingTimeout = 5 * time.Second

	logHelper.Info("Connect to database")

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logHelper.Info("Connect OKKK")

	queries := postgres.New(pool)

	d := &Data{
		db:      pool,
		queries: queries,
	}

	cleanup := func() {
		log.Info("closing the database connection pool")
		pool.Close()
	}

	return d, cleanup, nil
}
