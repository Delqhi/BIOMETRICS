package migration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMigrationDirectionString(t *testing.T) {
	tests := []struct {
		direction MigrationDirection
		expected  string
	}{
		{DirectionUp, "up"},
		{DirectionDown, "down"},
	}

	for _, tt := range tests {
		if string(tt.direction) != tt.expected {
			t.Errorf("Direction = %v, want %v", tt.direction, tt.expected)
		}
	}
}

func TestMigrationStatusString(t *testing.T) {
	tests := []struct {
		status   MigrationStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusApplied, "applied"},
		{StatusFailed, "failed"},
		{StatusReverted, "reverted"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("Status = %v, want %v", tt.status, tt.expected)
		}
	}
}

func TestDatabaseTypeString(t *testing.T) {
	tests := []struct {
		dbType   DatabaseType
		expected string
	}{
		{DBTypePostgres, "postgres"},
		{DBTypeSQLite, "sqlite"},
		{DBTypeMySQL, "mysql"},
	}

	for _, tt := range tests {
		if string(tt.dbType) != tt.expected {
			t.Errorf("DatabaseType = %v, want %v", tt.dbType, tt.expected)
		}
	}
}

func TestNewMigratorNilConfig(t *testing.T) {
	_, err := NewMigrator(nil)
	if err == nil {
		t.Error("NewMigrator() should fail with nil config")
	}
}

func TestNewMigratorNilDB(t *testing.T) {
	config := &Config{
		DB: nil,
	}
	_, err := NewMigrator(config)
	if err == nil {
		t.Error("NewMigrator() should fail with nil DB")
	}
}

func TestNewMigratorDefaultTableName(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
	}

	migrator, err := NewMigrator(config)
	if err != nil {
		t.Fatalf("NewMigrator() error = %v", err)
	}

	if migrator.config.TableName != "schema_migrations" {
		t.Errorf("TableName = %v, want schema_migrations", migrator.config.TableName)
	}
}

func TestMigratorEnsureMigrationTable(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		TableName:    "test_migrations",
	}

	migrator, err := NewMigrator(config)
	if err != nil {
		t.Fatalf("NewMigrator() error = %v", err)
	}
	_ = migrator

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query migration table: %v", err)
	}
}

func TestMigratorGetMigrationsEmpty(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	migrations := migrator.GetMigrations()
	if len(migrations) != 0 {
		t.Errorf("GetMigrations() = %d, want 0", len(migrations))
	}
}

func TestMigratorDiscoverMigrations(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)

	err = migrator.DiscoverMigrations()
	if err != nil {
		t.Fatalf("DiscoverMigrations() error = %v", err)
	}

	migrations := migrator.GetMigrations()
	if len(migrations) != 0 {
		t.Errorf("Expected 0 migrations, got %d", len(migrations))
	}
}

func TestMigratorDiscoverMigrationsWithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "1_create_users.down.sql"), []byte("DROP TABLE users;"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)

	err = migrator.DiscoverMigrations()
	if err != nil {
		t.Fatalf("DiscoverMigrations() error = %v", err)
	}

	migrations := migrator.GetMigrations()
	if len(migrations) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(migrations))
	}

	if migrations[0].Version != 1 {
		t.Errorf("First migration version = %d, want 1", migrations[0].Version)
	}

	if migrations[0].UpSQL == "" {
		t.Error("First migration should have UpSQL")
	}

	if migrations[0].DownSQL == "" {
		t.Error("First migration should have DownSQL")
	}
}

func TestMigratorGetPendingMigrations(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	pending := migrator.GetPendingMigrations()
	if len(pending) != 1 {
		t.Errorf("Expected 1 pending migration, got %d", len(pending))
	}
}

func TestMigratorCurrentVersionEmpty(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
	}

	migrator, _ := NewMigrator(config)

	version := migrator.CurrentVersion()
	if version != 0 {
		t.Errorf("CurrentVersion() = %d, want 0", version)
	}
}

func TestMigratorCurrentVersionWithMigrations(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
		Verbose:      true,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	results, err := migrator.Migrate(ctx, 0)
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 migrations, got %d", len(results))
	}

	version := migrator.CurrentVersion()
	if version != 2 {
		t.Errorf("CurrentVersion() = %d, want 2", version)
	}
}

func TestMigratorMigrateDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
		DryRun:       true,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	results, err := migrator.Migrate(ctx, 0)
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 migration result, got %d", len(results))
	}

	if results[0].Status != StatusPending {
		t.Errorf("Status = %v, want %v", results[0].Status, StatusPending)
	}

	version := migrator.CurrentVersion()
	if version != 0 {
		t.Errorf("CurrentVersion() = %d, want 0 (dry run)", version)
	}
}

func TestMigratorRollback(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "1_create_users.down.sql"), []byte("DROP TABLE users;"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
		Verbose:      true,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	_, err = migrator.Migrate(ctx, 0)
	if err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	versionBefore := migrator.CurrentVersion()
	if versionBefore != 1 {
		t.Fatalf("Expected version 1, got %d", versionBefore)
	}

	results, err := migrator.Rollback(ctx, 1)
	if err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 rollback result, got %d", len(results))
	}

	versionAfter := migrator.CurrentVersion()
	if versionAfter != 0 {
		t.Errorf("CurrentVersion() = %d, want 0 after rollback", versionAfter)
	}
}

func TestMigratorRollbackDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "1_create_users.down.sql"), []byte("DROP TABLE users;"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
		DryRun:       true,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	migrator.Migrate(ctx, 0)

	results, err := migrator.Rollback(ctx, 1)
	if err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 rollback result, got %d", len(results))
	}

	if results[0].Status != StatusReverted {
		t.Errorf("Status = %v, want %v", results[0].Status, StatusReverted)
	}

	version := migrator.CurrentVersion()
	if version != 1 {
		t.Errorf("CurrentVersion() = %d, want 1 (dry run)", version)
	}
}

func TestMigratorStatus(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	migrator.Migrate(context.Background(), 1)

	status, err := migrator.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if status[1] != StatusApplied {
		t.Errorf("status[1] = %v, want %v", status[1], StatusApplied)
	}

	if status[2] != StatusPending {
		t.Errorf("status[2] = %v, want %v", status[2], StatusPending)
	}
}

func TestMigratorForceVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
	}

	migrator, _ := NewMigrator(config)

	err = migrator.ForceVersion(5)
	if err != nil {
		t.Fatalf("ForceVersion() error = %v", err)
	}

	version := migrator.CurrentVersion()
	if version != 5 {
		t.Errorf("CurrentVersion() = %d, want 5", version)
	}
}

func TestMigratorValidateMigrations(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	errors, err := migrator.ValidateMigrations()
	if err != nil {
		t.Fatalf("ValidateMigrations() error = %v", err)
	}

	if len(errors) != 0 {
		t.Errorf("Expected no validation errors, got %d: %v", len(errors), errors)
	}
}

func TestMigratorPlanUp(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	plan, err := migrator.Plan(2)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	if plan.Direction != DirectionUp {
		t.Errorf("Direction = %v, want %v", plan.Direction, DirectionUp)
	}

	if len(plan.Migrations) != 2 {
		t.Errorf("Expected 2 migrations in plan, got %d", len(plan.Migrations))
	}

	if plan.ToVersion != 2 {
		t.Errorf("ToVersion = %d, want 2", plan.ToVersion)
	}
}

func TestMigratorPlanDown(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()
	migrator.Migrate(context.Background(), 2)

	plan, err := migrator.Plan(0)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	if plan.Direction != DirectionDown {
		t.Errorf("Direction = %v, want %v", plan.Direction, DirectionDown)
	}

	if len(plan.Migrations) != 2 {
		t.Errorf("Expected 2 migrations in plan, got %d", len(plan.Migrations))
	}
}

func TestMigratorGetAppliedMigrations(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	applied := migrator.GetAppliedMigrations()
	if len(applied) != 0 {
		t.Errorf("Expected 0 applied migrations, got %d", len(applied))
	}

	migrator.Migrate(context.Background(), 1)

	applied = migrator.GetAppliedMigrations()
	if len(applied) != 1 {
		t.Errorf("Expected 1 applied migration, got %d", len(applied))
	}

	if applied[0].Version != 1 {
		t.Errorf("Applied version = %d, want 1", applied[0].Version)
	}
}

func TestMigratorMultipleRollbacks(t *testing.T) {
	tmpDir := t.TempDir()

	for i := 1; i <= 5; i++ {
		err := os.WriteFile(
			filepath.Join(tmpDir, fmt.Sprintf("%d_test.up.sql", i)),
			[]byte(fmt.Sprintf("CREATE TABLE test_%d (id INTEGER PRIMARY KEY);", i)),
			0644,
		)
		if err != nil {
			t.Fatalf("Failed to write migration file: %v", err)
		}
		err = os.WriteFile(
			filepath.Join(tmpDir, fmt.Sprintf("%d_test.down.sql", i)),
			[]byte(fmt.Sprintf("DROP TABLE test_%d;", i)),
			0644,
		)
		if err != nil {
			t.Fatalf("Failed to write migration file: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	migrator.Migrate(ctx, 5)

	version := migrator.CurrentVersion()
	if version != 5 {
		t.Fatalf("Expected version 5, got %d", version)
	}

	migrator.Rollback(ctx, 2)
	version = migrator.CurrentVersion()
	if version != 3 {
		t.Errorf("CurrentVersion() = %d, want 3", version)
	}
}

func TestMigratorMigrateToSpecificVersion(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tmpDir, "1_create_users.up.sql"), []byte("CREATE TABLE users (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "2_create_posts.up.sql"), []byte("CREATE TABLE posts (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "3_create_comments.up.sql"), []byte("CREATE TABLE comments (id INTEGER PRIMARY KEY);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx := context.Background()
	migrator.Migrate(ctx, 2)

	version := migrator.CurrentVersion()
	if version != 2 {
		t.Errorf("CurrentVersion() = %d, want 2", version)
	}
}

func TestMigratorContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	for i := 1; i <= 3; i++ {
		err := os.WriteFile(
			filepath.Join(tmpDir, fmt.Sprintf("%d_test.up.sql", i)),
			[]byte(fmt.Sprintf("CREATE TABLE test_%d (id INTEGER PRIMARY KEY);", i)),
			0644,
		)
		if err != nil {
			t.Fatalf("Failed to write migration file: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	config := &Config{
		DB:           db,
		DatabaseType: DBTypeSQLite,
		Directory:    tmpDir,
	}

	migrator, _ := NewMigrator(config)
	migrator.DiscoverMigrations()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = migrator.Migrate(ctx, 3)
	if err == nil {
		t.Error("Migrate() should fail with cancelled context")
	}
}

func TestMigrationRecord(t *testing.T) {
	record := MigrationRecord{
		Version:       1,
		Name:          "test_migration",
		AppliedAt:     time.Now(),
		ExecutionTime: 100 * time.Millisecond,
		Success:       true,
	}

	if record.Version != 1 {
		t.Errorf("Version = %d, want 1", record.Version)
	}

	if record.Name != "test_migration" {
		t.Errorf("Name = %v, want test_migration", record.Name)
	}

	if !record.Success {
		t.Error("Success should be true")
	}
}

func TestMigration(t *testing.T) {
	migration := Migration{
		Version:       1,
		Name:          "create_users",
		UpSQL:         "CREATE TABLE users (id INTEGER PRIMARY KEY);",
		DownSQL:       "DROP TABLE users;",
		Direction:     DirectionUp,
		ExecutedAt:    time.Now(),
		Status:        StatusApplied,
		ExecutionTime: 50 * time.Millisecond,
	}

	if migration.Version != 1 {
		t.Errorf("Version = %d, want 1", migration.Version)
	}

	if migration.UpSQL == "" {
		t.Error("UpSQL should not be empty")
	}

	if migration.DownSQL == "" {
		t.Error("DownSQL should not be empty")
	}

	if migration.Status != StatusApplied {
		t.Errorf("Status = %v, want %v", migration.Status, StatusApplied)
	}
}
