package database

import (
	"context"
	"fmt"
	"time"

	"biometrics/internal/config"
	"biometrics/pkg/models"

	"biometrics/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	pool   *pgxpool.Pool
	logger utils.Logger
}

func NewPostgres(databaseURL string, logger utils.Logger) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully",
		"max_conns", config.MaxConns,
		"min_conns", config.MinConns,
	)

	return &Postgres{
		pool:   pool,
		logger: logger,
	}, nil
}

func NewPostgresFromConfig(cfg config.DatabaseConfig, logger utils.Logger) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	return NewPostgres(dsn, logger)
}

func (p *Postgres) Close() {
	p.pool.Close()
	p.logger.Info("Database connection closed")
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) Query(ctx context.Context, sql string, args ...interface{}) (stdlib.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *Postgres) QueryRow(ctx context.Context, sql string, args ...interface{}) stdlib.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *Postgres) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	result, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (p *Postgres) QueryWithResult(ctx context.Context, sql string, args ...interface{}) (int, []map[string]interface{}, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	var results []map[string]interface{}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return 0, nil, err
		}

		row := make(map[string]interface{})
		for i, field := range fields {
			row[string(field.Name)] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return len(fields), results, nil
}

func (p *Postgres) Acquire(ctx context.Context) (stdlib.Conn, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	return stdlib.Conn{Conn: conn.Conn()}, nil
}

func (p *Postgres) Stats() *pgxpool.Stat {
	return p.pool.Stat()
}

func (p *Postgres) Begin(ctx context.Context) (stdlib.Tx, error) {
	return p.pool.Begin(ctx)
}

func (p *Postgres) BeginTx(ctx context.Context, txOptions TxOptions) (stdlib.Tx, error) {
	return p.pool.BeginTx(ctx, txOptions)
}

type TxOptions struct {
	IsoLevel       string
	AccessMode     string
	DeferrableMode string
}

func (p *Postgres) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	sql := `SELECT id, email, password_hash, first_name, last_name, phone, role, status, 
			email_verified, two_factor_enabled, provider, provider_id, avatar, bio, roles, 
			permissions, tenant_id, metadata, last_login_at, failed_login_attempts, 
			locked_until, created_at, updated_at, deleted_at
			FROM users WHERE id = $1 AND deleted_at IS NULL`

	row := p.QueryRow(ctx, sql, id)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.EmailVerified, &user.TwoFactorEnabled,
		&user.Provider, &user.ProviderID, &user.Avatar, &user.Bio, &user.Roles,
		&user.Permissions, &user.TenantID, &user.Metadata, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, email, password_hash, first_name, last_name, phone, role, status, 
			email_verified, two_factor_enabled, provider, provider_id, avatar, bio, roles, 
			permissions, tenant_id, metadata, last_login_at, failed_login_attempts, 
			locked_until, created_at, updated_at, deleted_at
			FROM users WHERE email = $1 AND deleted_at IS NULL`

	row := p.QueryRow(ctx, sql, email)

	var user models.User
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.Status, &user.EmailVerified, &user.TwoFactorEnabled,
		&user.Provider, &user.ProviderID, &user.Avatar, &user.Bio, &user.Roles,
		&user.Permissions, &user.TenantID, &user.Metadata, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *models.User) error {
	sql := `INSERT INTO users (id, email, password_hash, first_name, last_name, phone, role, status,
			email_verified, two_factor_enabled, provider, provider_id, avatar, bio, roles, 
			permissions, tenant_id, metadata, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	_, err := p.Exec(ctx, sql,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Role, user.Status, user.EmailVerified, user.TwoFactorEnabled,
		user.Provider, user.ProviderID, user.Avatar, user.Bio, user.Roles,
		user.Permissions, user.TenantID, user.Metadata, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (p *Postgres) UpdateUser(ctx context.Context, user *models.User) error {
	sql := `UPDATE users SET 
			email = $1, first_name = $2, last_name = $3, phone = $4, role = $5, status = $6,
			email_verified = $7, two_factor_enabled = $8, avatar = $9, bio = $10, roles = $11,
			permissions = $12, tenant_id = $13, metadata = $14, last_login_at = $15,
			failed_login_attempts = $16, locked_until = $17, updated_at = $18
			WHERE id = $19`

	_, err := p.Exec(ctx, sql,
		user.Email, user.FirstName, user.LastName, user.Phone, user.Role, user.Status,
		user.EmailVerified, user.TwoFactorEnabled, user.Avatar, user.Bio, user.Roles,
		user.Permissions, user.TenantID, user.Metadata, user.LastLoginAt,
		user.FailedLoginAttempts, user.LockedUntil, user.UpdatedAt, user.ID,
	)
	return err
}

func (p *Postgres) DeleteUser(ctx context.Context, id string) error {
	sql := `UPDATE users SET deleted_at = $1, updated_at = $1 WHERE id = $2`
	_, err := p.Exec(ctx, sql, time.Now(), id)
	return err
}

func (p *Postgres) ListUsers(ctx context.Context, filter models.UserFilter) ([]*models.User, int64, error) {
	baseSQL := `FROM users WHERE deleted_at IS NULL`
	countSQL := `SELECT COUNT(*) ` + baseSQL
	selectSQL := `SELECT id, email, password_hash, first_name, last_name, phone, role, status, 
			email_verified, two_factor_enabled, provider, provider_id, avatar, bio, roles, 
			permissions, tenant_id, metadata, last_login_at, failed_login_attempts, 
			locked_until, created_at, updated_at, deleted_at ` + baseSQL

	var args []interface{}
	argIndex := 1

	if filter.Email != "" {
		selectSQL += fmt.Sprintf(" AND email = $%d", argIndex)
		countSQL += fmt.Sprintf(" AND email = $%d", argIndex)
		args = append(args, filter.Email)
		argIndex++
	}

	if filter.Role != "" {
		selectSQL += fmt.Sprintf(" AND role = $%d", argIndex)
		countSQL += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, filter.Role)
		argIndex++
	}

	if filter.Status != "" {
		selectSQL += fmt.Sprintf(" AND status = $%d", argIndex)
		countSQL += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	var total int64
	p.QueryRow(ctx, countSQL, args...).Scan(&total)

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	selectSQL += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.PageSize, offset)

	rows, err := p.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.Phone, &user.Role, &user.Status, &user.EmailVerified, &user.TwoFactorEnabled,
			&user.Provider, &user.ProviderID, &user.Avatar, &user.Bio, &user.Roles,
			&user.Permissions, &user.TenantID, &user.Metadata, &user.LastLoginAt,
			&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	return users, total, nil
}
