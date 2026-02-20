// Package migration provides database migration capabilities with support for
// version tracking, up/down migrations, rollbacks, and multiple database backends.
package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MigrationDirection string

const (
	DirectionUp   MigrationDirection = "up"
	DirectionDown MigrationDirection = "down"
)

type MigrationStatus string

const (
	StatusPending  MigrationStatus = "pending"
	StatusApplied  MigrationStatus = "applied"
	StatusFailed   MigrationStatus = "failed"
	StatusReverted MigrationStatus = "reverted"
)

type DatabaseType string

const (
	DBTypePostgres DatabaseType = "postgres"
	DBTypeSQLite   DatabaseType = "sqlite"
	DBTypeMySQL    DatabaseType = "mysql"
)

type Migration struct {
	Version       int
	Name          string
	UpSQL         string
	DownSQL       string
	Direction     MigrationDirection
	ExecutedAt    time.Time
	ExecutionTime time.Duration
	Status        MigrationStatus
	Error         string
}

type MigrationFile struct {
	Version     int
	Name        string
	Path        string
	UpContent   string
	DownContent string
}

type MigrationRecord struct {
	Version       int
	Name          string
	AppliedAt     time.Time
	ExecutionTime time.Duration
	Success       bool
	Error         string
}

type Config struct {
	TableName    string
	Directory    string
	DatabaseType DatabaseType
	DB           *sql.DB
	AutoDump     bool
	DryRun       bool
	Verbose      bool
}

type Migrator struct {
	config     *Config
	mu         sync.RWMutex
	migrations []Migration
	applied    map[int]MigrationRecord
}

func NewMigrator(config *Config) (*Migrator, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if config.DB == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}
	if config.TableName == "" {
		config.TableName = "schema_migrations"
	}
	if config.Directory == "" {
		config.Directory = "./migrations"
	}

	m := &Migrator{
		config:     config,
		migrations: make([]Migration, 0),
		applied:    make(map[int]MigrationRecord),
	}

	if err := m.ensureMigrationTable(); err != nil {
		return nil, fmt.Errorf("failed to ensure migration table: %w", err)
	}

	if err := m.loadAppliedMigrations(); err != nil {
		return nil, fmt.Errorf("failed to load applied migrations: %w", err)
	}

	return m, nil
}

func (m *Migrator) ensureMigrationTable() error {
	var createSQL string

	switch m.config.DatabaseType {
	case DBTypePostgres:
		createSQL = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				version INTEGER PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				applied_at TIMESTAMP NOT NULL DEFAULT NOW(),
				execution_time_ms INTEGER,
				success BOOLEAN NOT NULL DEFAULT true,
				error TEXT
			)
		`, m.config.TableName)
	case DBTypeSQLite:
		createSQL = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				version INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				applied_at TEXT NOT NULL,
				execution_time_ms INTEGER,
				success INTEGER NOT NULL DEFAULT 1,
				error TEXT
			)
		`, m.config.TableName)
	case DBTypeMySQL:
		createSQL = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				version INT PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				execution_time_ms INT,
				success TINYINT(1) NOT NULL DEFAULT 1,
				error TEXT
			)
		`, m.config.TableName)
	default:
		return fmt.Errorf("unsupported database type: %s", m.config.DatabaseType)
	}

	_, err := m.config.DB.Exec(createSQL)
	return err
}

func (m *Migrator) loadAppliedMigrations() error {
	query := fmt.Sprintf("SELECT version, name, applied_at, execution_time_ms, success, error FROM %s", m.config.TableName)

	rows, err := m.config.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record MigrationRecord
		var execTime sql.NullInt64
		var success int

		err := rows.Scan(&record.Version, &record.Name, &record.AppliedAt, &execTime, &success, &record.Error)
		if err != nil {
			return err
		}

		record.Success = success == 1
		if execTime.Valid {
			record.ExecutionTime = time.Duration(execTime.Int64) * time.Millisecond
		}

		m.applied[record.Version] = record
	}

	return rows.Err()
}

func (m *Migrator) DiscoverMigrations() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	files, err := ioutil.ReadDir(m.config.Directory)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	var migrationFiles []MigrationFile

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		versionStr := strings.TrimSuffix(parts[0], ".up")
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			continue
		}

		migrationName := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")

		upPath := filepath.Join(m.config.Directory, fmt.Sprintf("%d_%s.up.sql", version, migrationName))
		downPath := filepath.Join(m.config.Directory, fmt.Sprintf("%d_%s.down.sql", version, migrationName))

		upContent := []byte{}
		if _, err := os.Stat(upPath); err == nil {
			upContent, err = ioutil.ReadFile(upPath)
			if err != nil {
				return fmt.Errorf("failed to read up migration %s: %w", upPath, err)
			}
		}

		downContent := []byte{}
		if _, err := os.Stat(downPath); err == nil {
			downContent, err = ioutil.ReadFile(downPath)
			if err != nil {
				return fmt.Errorf("failed to read down migration %s: %w", downPath, err)
			}
		}

		migrationFiles = append(migrationFiles, MigrationFile{
			Version:     version,
			Name:        migrationName,
			Path:        filepath.Join(m.config.Directory, name),
			UpContent:   string(upContent),
			DownContent: string(downContent),
		})
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Version < migrationFiles[j].Version
	})

	m.migrations = make([]Migration, len(migrationFiles))
	for i, mf := range migrationFiles {
		m.migrations[i] = Migration{
			Version: mf.Version,
			Name:    mf.Name,
			UpSQL:   mf.UpContent,
			DownSQL: mf.DownContent,
		}
	}

	return nil
}

func (m *Migrator) GetMigrations() []Migration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]Migration, len(m.migrations))
	copy(result, m.migrations)
	return result
}

func (m *Migrator) GetPendingMigrations() []Migration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var pending []Migration
	for _, migration := range m.migrations {
		if _, applied := m.applied[migration.Version]; !applied {
			pending = append(pending, migration)
		}
	}

	return pending
}

func (m *Migrator) GetAppliedMigrations() []MigrationRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records := make([]MigrationRecord, 0, len(m.applied))
	for _, record := range m.applied {
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Version < records[j].Version
	})

	return records
}

func (m *Migrator) Migrate(ctx context.Context, targetVersion int) ([]Migration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pending := m.getPendingMigrationsLocked()
	if len(pending) == 0 {
		return nil, nil
	}

	var toApply []Migration
	for _, migration := range pending {
		if targetVersion > 0 && migration.Version > targetVersion {
			break
		}
		toApply = append(toApply, migration)
	}

	if len(toApply) == 0 {
		return nil, nil
	}

	results := make([]Migration, 0, len(toApply))

	for _, migration := range toApply {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		if m.config.DryRun {
			migration.Status = StatusPending
			results = append(results, migration)
			continue
		}

		result, err := m.executeMigrationLocked(ctx, migration, DirectionUp)
		results = append(results, result)

		if err != nil {
			return results, fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}
	}

	return results, nil
}

func (m *Migrator) Rollback(ctx context.Context, steps int) ([]Migration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if steps <= 0 {
		return nil, nil
	}

	applied := m.getAppliedMigrationsLocked()
	if len(applied) == 0 {
		return nil, nil
	}

	toRollback := applied
	if len(toRollback) > steps {
		toRollback = toRollback[len(toRollback)-steps:]
	}

	results := make([]Migration, 0, len(toRollback))

	for i := len(toRollback) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		record := toRollback[i]
		migration := Migration{
			Version: record.Version,
			Name:    record.Name,
			DownSQL: "",
		}

		for _, m := range m.migrations {
			if m.Version == migration.Version {
				migration.DownSQL = m.DownSQL
				break
			}
		}

		if migration.DownSQL == "" {
			continue
		}

		if m.config.DryRun {
			migration.Status = StatusReverted
			results = append(results, migration)
			continue
		}

		result, err := m.executeMigrationLocked(ctx, migration, DirectionDown)
		results = append(results, result)

		if err != nil {
			return results, fmt.Errorf("rollback %d failed: %w", migration.Version, err)
		}
	}

	return results, nil
}

func (m *Migrator) executeMigrationLocked(ctx context.Context, migration Migration, direction MigrationDirection) (Migration, error) {
	migration.Direction = direction
	migration.ExecutedAt = time.Now()

	sqlToExecute := migration.UpSQL
	if direction == DirectionDown {
		sqlToExecute = migration.DownSQL
	}

	if sqlToExecute == "" {
		migration.Status = StatusFailed
		migration.Error = "no SQL to execute"
		return migration, fmt.Errorf("no %s SQL for migration %d", direction, migration.Version)
	}

	startTime := time.Now()

	tx, err := m.config.DB.BeginTx(ctx, nil)
	if err != nil {
		migration.Status = StatusFailed
		migration.Error = err.Error()
		return migration, fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, sqlToExecute)
	if err != nil {
		tx.Rollback()
		migration.Status = StatusFailed
		migration.Error = err.Error()
		migration.ExecutionTime = time.Since(startTime)

		m.recordMigrationLocked(migration)
		return migration, err
	}

	if err := tx.Commit(); err != nil {
		migration.Status = StatusFailed
		migration.Error = err.Error()
		migration.ExecutionTime = time.Since(startTime)

		m.recordMigrationLocked(migration)
		return migration, err
	}

	migration.ExecutionTime = time.Since(startTime)
	migration.Status = StatusApplied

	m.recordMigrationLocked(migration)

	if direction == DirectionUp {
		m.applied[migration.Version] = MigrationRecord{
			Version:       migration.Version,
			Name:          migration.Name,
			AppliedAt:     migration.ExecutedAt,
			ExecutionTime: migration.ExecutionTime,
			Success:       true,
		}
	} else {
		delete(m.applied, migration.Version)
	}

	if m.config.Verbose {
		fmt.Printf("Applied migration %d (%s) in %v\n", migration.Version, migration.Name, migration.ExecutionTime)
	}

	return migration, nil
}

func (m *Migrator) recordMigrationLocked(migration Migration) error {
	success := migration.Status == StatusApplied

	var execTimeMs int64
	if migration.ExecutionTime > 0 {
		execTimeMs = migration.ExecutionTime.Milliseconds()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (version, name, applied_at, execution_time_ms, success, error)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (version) DO UPDATE SET
			applied_at = $3,
			execution_time_ms = $4,
			success = $5,
			error = $6
	`, m.config.TableName)

	_, err := m.config.DB.Exec(query,
		migration.Version,
		migration.Name,
		migration.ExecutedAt,
		execTimeMs,
		success,
		migration.Error,
	)

	return err
}

func (m *Migrator) getPendingMigrationsLocked() []Migration {
	var pending []Migration
	for _, migration := range m.migrations {
		if _, applied := m.applied[migration.Version]; !applied {
			pending = append(pending, migration)
		}
	}
	return pending
}

func (m *Migrator) getAppliedMigrationsLocked() []MigrationRecord {
	records := make([]MigrationRecord, 0, len(m.applied))
	for _, record := range m.applied {
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Version < records[j].Version
	})

	return records
}

func (m *Migrator) CurrentVersion() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.applied) == 0 {
		return 0
	}

	maxVersion := 0
	for version := range m.applied {
		if version > maxVersion {
			maxVersion = version
		}
	}

	return maxVersion
}

func (m *Migrator) Status() (map[int]MigrationStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[int]MigrationStatus)

	for _, migration := range m.migrations {
		if record, applied := m.applied[migration.Version]; applied {
			if record.Success {
				status[migration.Version] = StatusApplied
			} else {
				status[migration.Version] = StatusFailed
			}
		} else {
			status[migration.Version] = StatusPending
		}
	}

	return status, nil
}

func (m *Migrator) ForceVersion(version int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var migration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == version {
			migration = &m.migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration version %d not found", version)
	}

	record := MigrationRecord{
		Version:   migration.Version,
		Name:      migration.Name,
		AppliedAt: time.Now(),
		Success:   true,
	}

	m.applied[version] = record

	query := fmt.Sprintf(`
		INSERT INTO %s (version, name, applied_at, execution_time_ms, success, error)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (version) DO UPDATE SET
			applied_at = $3,
			success = $5,
			error = $6
	`, m.config.TableName)

	_, err := m.config.DB.Exec(query, record.Version, record.Name, record.AppliedAt, 0, true, nil)
	return err
}

func (m *Migrator) ValidateMigrations() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []string

	versionPattern := regexp.MustCompile(`^\d+_`)

	for _, migration := range m.migrations {
		if migration.Version <= 0 {
			errors = append(errors, fmt.Sprintf("migration %s has invalid version %d", migration.Name, migration.Version))
		}

		if !versionPattern.MatchString(migration.Name) {
			errors = append(errors, fmt.Sprintf("migration %s does not match naming convention (version_name.sql)", migration.Name))
		}

		if migration.UpSQL == "" && migration.DownSQL == "" {
			errors = append(errors, fmt.Sprintf("migration %d has no up or down SQL", migration.Version))
		}
	}

	versions := make(map[int]bool)
	for _, migration := range m.migrations {
		if versions[migration.Version] {
			errors = append(errors, fmt.Sprintf("duplicate migration version %d", migration.Version))
		}
		versions[migration.Version] = true
	}

	return errors, nil
}

type MigrationPlan struct {
	FromVersion int
	ToVersion   int
	Direction   MigrationDirection
	Migrations  []Migration
}

func (m *Migrator) Plan(targetVersion int) (*MigrationPlan, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	currentVersion := m.CurrentVersion()

	plan := &MigrationPlan{
		FromVersion: currentVersion,
		ToVersion:   targetVersion,
	}

	if targetVersion > currentVersion {
		plan.Direction = DirectionUp
		for _, migration := range m.migrations {
			if migration.Version > currentVersion && (targetVersion == 0 || migration.Version <= targetVersion) {
				if _, applied := m.applied[migration.Version]; !applied {
					plan.Migrations = append(plan.Migrations, migration)
				}
			}
		}
	} else if targetVersion < currentVersion {
		plan.Direction = DirectionDown
		for i := len(m.migrations) - 1; i >= 0; i-- {
			migration := m.migrations[i]
			if migration.Version <= currentVersion && migration.Version > targetVersion {
				if _, applied := m.applied[migration.Version]; applied {
					plan.Migrations = append(plan.Migrations, migration)
				}
			}
		}
	}

	return plan, nil
}
