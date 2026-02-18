package database

import (
	"context"
	"fmt"
	"time"

	"biometrics/pkg/utils"
)

type Migration struct {
	Version    string
	Name       string
	ExecutedAt time.Time
}

type MigrationManager struct {
	db     *Postgres
	logger utils.Logger
}

func NewMigrationManager(db *Postgres, logger utils.Logger) *MigrationManager {
	return &MigrationManager{
		db:     db,
		logger: logger,
	}
}

func (m *MigrationManager) RunMigrations(ctx context.Context) error {
	m.logger.Info("Starting database migrations...")

	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := []MigrationScript{
		{
			Version: "001",
			Name:    "create_users_table",
			Up:      UsersTableUp,
			Down:    UsersTableDown,
		},
		{
			Version: "002",
			Name:    "create_sessions_table",
			Up:      SessionsTableUp,
			Down:    SessionsTableDown,
		},
		{
			Version: "003",
			Name:    "create_content_table",
			Up:      ContentTableUp,
			Down:    ContentTableDown,
		},
		{
			Version: "004",
			Name:    "create_integrations_table",
			Up:      IntegrationsTableUp,
			Down:    IntegrationsTableDown,
		},
		{
			Version: "005",
			Name:    "create_workflows_table",
			Up:      WorkflowsTableUp,
			Down:    WorkflowsTableDown,
		},
		{
			Version: "006",
			Name:    "create_biometrics_table",
			Up:      BiometricsTableUp,
			Down:    BiometricsTableDown,
		},
		{
			Version: "007",
			Name:    "create_audit_logs_table",
			Up:      AuditLogsTableUp,
			Down:    AuditLogsTableDown,
		},
		{
			Version: "008",
			Name:    "create_tokens_table",
			Up:      TokensTableUp,
			Down:    TokensTableDown,
		},
		{
			Version: "009",
			Name:    "create_indexes",
			Up:      IndexesUp,
			Down:    IndexesDown,
		},
	}

	for _, migration := range migrations {
		if err := m.runMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
		}
	}

	m.logger.Info("Database migrations completed successfully")
	return nil
}

func (m *MigrationManager) createMigrationsTable(ctx context.Context) error {
	sql := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err := m.db.Exec(ctx, sql)
	return err
}

func (m *MigrationManager) runMigration(ctx context.Context, script MigrationScript) error {
	var count int
	err := m.db.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", script.Version).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration: %w", err)
	}

	if count > 0 {
		m.logger.Info("Migration already executed", "version", script.Version, "name", script.Name)
		return nil
	}

	m.logger.Info("Running migration", "version", script.Version, "name", script.Name)

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.Exec(ctx, script.Up); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute up migration: %w", err)
	}

	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, name, executed_at) VALUES ($1, $2, $3)",
		script.Version, script.Name, time.Now()); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	m.logger.Info("Migration completed", "version", script.Version, "name", script.Name)
	return nil
}

type MigrationScript struct {
	Version string
	Name    string
	Up      string
	Down    string
}

const UsersTableUp = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    role VARCHAR(50) DEFAULT 'user',
    status VARCHAR(50) DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    provider VARCHAR(50) DEFAULT 'local',
    provider_id VARCHAR(255),
    avatar TEXT,
    bio TEXT,
    roles TEXT[] DEFAULT ARRAY['user'],
    permissions TEXT[] DEFAULT ARRAY['read'],
    tenant_id UUID,
    metadata JSONB DEFAULT '{}',
    last_login_at TIMESTAMP WITH TIME ZONE,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
`

const UsersTableDown = `
DROP TABLE IF EXISTS users;
`

const SessionsTableUp = `
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    refresh_token_hash VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_id VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
`

const SessionsTableDown = `DROP TABLE IF EXISTS sessions;`

const ContentTableUp = `
CREATE TABLE IF NOT EXISTS content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    body TEXT,
    excerpt TEXT,
    type VARCHAR(50) DEFAULT 'article',
    status VARCHAR(50) DEFAULT 'draft',
    author_id UUID NOT NULL REFERENCES users(id),
    featured_image TEXT,
    meta_title VARCHAR(200),
    meta_description VARCHAR(500),
    tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    categories TEXT[] DEFAULT ARRAY[]::TEXT[],
    published_at TIMESTAMP WITH TIME ZONE,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    language VARCHAR(10) DEFAULT 'en',
    version INTEGER DEFAULT 1,
    parent_id UUID REFERENCES content(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_content_tenant_id ON content(tenant_id);
CREATE INDEX idx_content_author_id ON content(author_id);
CREATE INDEX idx_content_status ON content(status);
CREATE INDEX idx_content_type ON content(type);
CREATE INDEX idx_content_slug ON content(slug);
`

const ContentTableDown = `DROP TABLE IF EXISTS content;`

const IntegrationsTableUp = `
CREATE TABLE IF NOT EXISTS integrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    logo TEXT,
    website VARCHAR(500),
    owner_id UUID NOT NULL REFERENCES users(id),
    api_key_hash VARCHAR(255),
    api_secret_hash VARCHAR(255),
    access_token_hash VARCHAR(255),
    refresh_token_hash VARCHAR(255),
    webhook_url VARCHAR(500),
    webhook_secret_hash VARCHAR(255),
    scopes TEXT[] DEFAULT ARRAY[]::TEXT[],
    settings JSONB DEFAULT '{}',
    credentials_encrypted TEXT,
    rate_limit INTEGER DEFAULT 1000,
    rate_remaining INTEGER DEFAULT 1000,
    rate_reset_at TIMESTAMP WITH TIME ZONE,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    next_sync_at TIMESTAMP WITH TIME ZONE,
    sync_interval INTEGER DEFAULT 3600,
    failure_count INTEGER DEFAULT 0,
    last_error TEXT,
    metadata JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_integrations_tenant_id ON integrations(tenant_id);
CREATE INDEX idx_integrations_owner_id ON integrations(owner_id);
CREATE INDEX idx_integrations_type ON integrations(type);
CREATE INDEX idx_integrations_provider ON integrations(provider);
CREATE INDEX idx_integrations_status ON integrations(status);
`

const IntegrationsTableDown = `DROP TABLE IF EXISTS integrations;`

const WorkflowsTableUp = `
CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    definition JSONB,
    trigger_type VARCHAR(50) DEFAULT 'manual',
    trigger_config JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'draft',
    owner_id UUID NOT NULL REFERENCES users(id),
    version INTEGER DEFAULT 1,
    nodes_count INTEGER DEFAULT 0,
    edges_count INTEGER DEFAULT 0,
    last_run_at TIMESTAMP WITH TIME ZONE,
    next_run_at TIMESTAMP WITH TIME ZONE,
    run_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    failure_count INTEGER DEFAULT 0,
    avg_duration INTEGER DEFAULT 0,
    schedule VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    is_published BOOLEAN DEFAULT FALSE,
    execution_mode VARCHAR(50) DEFAULT 'sync',
    timeout INTEGER DEFAULT 300,
    retry_policy JSONB DEFAULT '{"max_retries": 3, "retry_delay": 1000}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_workflows_tenant_id ON workflows(tenant_id);
CREATE INDEX idx_workflows_owner_id ON workflows(owner_id);
CREATE INDEX idx_workflows_status ON workflows(status);
CREATE INDEX idx_workflows_trigger_type ON workflows(trigger_type);
`

const WorkflowsTableDown = `DROP TABLE IF EXISTS workflows;`

const BiometricsTableUp = `
CREATE TABLE IF NOT EXISTS biometrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID,
    type VARCHAR(50) NOT NULL,
    sub_type VARCHAR(50),
    label VARCHAR(100),
    template_encrypted TEXT,
    algorithm VARCHAR(50) DEFAULT 'fido2',
    device_id VARCHAR(255),
    device_info TEXT,
    public_key TEXT,
    credential_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'enrolled',
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    failure_count INTEGER DEFAULT 0,
    max_failures INTEGER DEFAULT 5,
    verify_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    quality_score INTEGER DEFAULT 100,
    metadata JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    is_primary BOOLEAN DEFAULT FALSE,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS biometric_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    biometric_id UUID NOT NULL REFERENCES biometrics(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    tenant_id UUID,
    challenge VARCHAR(255),
    challenge_hash VARCHAR(255),
    client_data TEXT,
    auth_data TEXT,
    signature_encrypted TEXT,
    result VARCHAR(50) NOT NULL,
    score INTEGER DEFAULT 0,
    failure_reason TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_info TEXT,
    location VARCHAR(255),
    latency INTEGER DEFAULT 0,
    verified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_biometrics_user_id ON biometrics(user_id);
CREATE INDEX idx_biometrics_tenant_id ON biometrics(tenant_id);
CREATE INDEX idx_biometrics_type ON biometrics(type);
CREATE INDEX idx_biometrics_status ON biometrics(status);
CREATE INDEX idx_biometric_verifications_biometric_id ON biometric_verifications(biometric_id);
CREATE INDEX idx_biometric_verifications_user_id ON biometric_verifications(user_id);
CREATE INDEX idx_biometric_verifications_result ON biometric_verifications(result);
`

const BiometricsTableDown = `
DROP TABLE IF EXISTS biometric_verifications;
DROP TABLE IF EXISTS biometrics;
`

const AuditLogsTableUp = `
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    user_id UUID,
    actor_id UUID,
    action VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    description TEXT,
    old_values JSONB,
    new_values JSONB,
    changes JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    location VARCHAR(255),
    session_id VARCHAR(255),
    request_id VARCHAR(255),
    correlation_id VARCHAR(255),
    level VARCHAR(50) DEFAULT 'info',
    severity VARCHAR(50) DEFAULT 'low',
    status VARCHAR(50) DEFAULT 'success',
    error_message TEXT,
    duration INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT ARRAY[]::TEXT[],
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_category ON audit_logs(category);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_level ON audit_logs(level);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
`

const AuditLogsTableDown = `DROP TABLE IF EXISTS audit_logs;`

const TokensTableUp = `
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID,
    type VARCHAR(50) NOT NULL,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    token_prefix VARCHAR(20),
    refresh_token_hash VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_id VARCHAR(255),
    device_info TEXT,
    location VARCHAR(255),
    scopes TEXT[] DEFAULT ARRAY[]::TEXT[],
    status VARCHAR(50) DEFAULT 'active',
    issued_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    revoke_reason VARCHAR(255),
    use_count INTEGER DEFAULT 0,
    max_uses INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_tenant_id ON tokens(tenant_id);
CREATE INDEX idx_tokens_token_hash ON tokens(token_hash);
CREATE INDEX idx_tokens_type ON tokens(type);
CREATE INDEX idx_tokens_status ON tokens(status);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
`

const TokensTableDown = `DROP TABLE IF EXISTS tokens;`

const IndexesUp = `
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users(LOWER(email));
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_created_at ON content(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_published_at ON content(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_integrations_created_at ON integrations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_workflows_created_at ON workflows(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_created ON audit_logs(tenant_id, created_at DESC);
`

const IndexesDown = `
DROP INDEX IF EXISTS idx_users_email_lower;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_content_created_at;
DROP INDEX IF EXISTS idx_content_published_at;
DROP INDEX IF EXISTS idx_integrations_created_at;
DROP INDEX IF EXISTS idx_workflows_created_at;
DROP INDEX IF EXISTS idx_audit_logs_tenant_created;
`
