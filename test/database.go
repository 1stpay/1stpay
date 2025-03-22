package test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/1stpay/1stpay/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AdminDatabase struct {
	db *sql.DB
}

func NewAdminDatabase(connStr string) (*AdminDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to admin DB: %w", err)
	}
	return &AdminDatabase{db: db}, nil
}

func (adminDB *AdminDatabase) CreateTestDatabase(ctx context.Context) (string, error) {
	testDBName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())
	createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", testDBName)
	if _, err := adminDB.db.ExecContext(ctx, createQuery); err != nil {
		return "", fmt.Errorf("error creating test DB: %w", err)
	}
	return testDBName, nil
}

type TestDatabase struct {
	TestDBName  string
	SQLDB       *sql.DB
	GormDB      *gorm.DB
	AdminDB     *AdminDatabase
	TestConnStr string
}

func NewTestPostgresDatabase(ctx context.Context, env *config.Env, rootPath string) (*TestDatabase, error) {

	adminConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)
	adminDB, err := NewAdminDatabase(adminConnStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to admin database: %w", err)
	}

	testDBName, err := adminDB.CreateTestDatabase(ctx)
	fmt.Printf("Test db name %v", testDBName)
	if err != nil {
		return nil, fmt.Errorf("error creating test database: %w", err)
	}

	testConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPass, env.DBHost, env.DBPort, testDBName)

	sqlDB, err := sql.Open("postgres", testConnStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to test database: %w", err)
	}

	gormDB, err := gorm.Open(postgres.Open(testConnStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating GORM connection: %w", err)
	}

	testDatabase := &TestDatabase{
		TestDBName:  testDBName,
		SQLDB:       sqlDB,
		GormDB:      gormDB,
		AdminDB:     adminDB,
		TestConnStr: testConnStr,
	}

	migrationsPath := fmt.Sprintf("%vmigrations", rootPath)
	if err := testDatabase.ApplyMigrations(ctx, migrationsPath); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	return testDatabase, nil
}

func (td *TestDatabase) ApplyMigrations(ctx context.Context, migrationsDir string) error {
	goose.SetDialect("postgres")
	if err := goose.Up(td.SQLDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func (td *TestDatabase) Cleanup(ctx context.Context) error {

	if err := td.SQLDB.Close(); err != nil {
		return fmt.Errorf("error closing test DB connection: %w", err)
	}
	terminateQuery := fmt.Sprintf(`
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = '%s' AND pid <> pg_backend_pid();
	`, td.TestDBName)
	if _, err := td.AdminDB.db.ExecContext(ctx, terminateQuery); err != nil {
		return fmt.Errorf("error terminating connections to test database: %w", err)
	}

	dropQuery := fmt.Sprintf("DROP DATABASE \"%s\"", td.TestDBName)
	if _, err := td.AdminDB.db.ExecContext(ctx, dropQuery); err != nil {
		return fmt.Errorf("error dropping test database: %w", err)
	}

	if err := td.AdminDB.db.Close(); err != nil {
		return fmt.Errorf("error closing admin DB connection: %w", err)
	}
	return nil
}
