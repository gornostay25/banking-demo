package testhelpers

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"banking/internal/database"
	"banking/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
)

import _ "unsafe"

//go:linkname dbDatabase banking/internal/database.database
var dbDatabase string

//go:linkname dbPassword banking/internal/database.password
var dbPassword string

//go:linkname dbUsername banking/internal/database.username
var dbUsername string

//go:linkname dbHost banking/internal/database.host
var dbHost string

//go:linkname dbPort banking/internal/database.port
var dbPort string

//go:linkname dbSchema banking/internal/database.schema
var dbSchema string

const (
	initMigrationFile     = "20260129170243_init.sql"
	testDataMigrationFile = "20260130145301_test_data.sql"
)

var (
	setupOnce   sync.Once
	setupErr    error
	testDB      *bun.DB
	cleanupFunc func() error
)

// checkDockerAvailable checks if Docker is available and accessible.
func checkDockerAvailable() error {
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker is not available or not accessible: %w. Please ensure Docker is running and you have permissions to access it", err)
	}
	return nil
}

// SetupTestDB initializes a postgres container and runs migrations once per package.
func SetupTestDB() (*bun.DB, func() error, error) {
	setupOnce.Do(func() {
		gin.SetMode(gin.TestMode)
		_ = os.Setenv("JWT_SECRET_KEY", "test-secret-key")

		// Check Docker availability before attempting to start container
		if err := checkDockerAvailable(); err != nil {
			setupErr = err
			return
		}

		ctx := context.Background()
		dbName := "banking"
		dbPwd := "password"
		dbUser := "user"

		dbContainer, err := postgres.Run(
			ctx,
			"postgres:latest",
			postgres.WithDatabase(dbName),
			postgres.WithUsername(dbUser),
			postgres.WithPassword(dbPwd),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		if err != nil {
			setupErr = fmt.Errorf("failed to start postgres container: %w", err)
			return
		}

		host, err := dbContainer.Host(ctx)
		if err != nil {
			_ = dbContainer.Terminate(ctx)
			setupErr = fmt.Errorf("failed to get container host: %w", err)
			return
		}

		port, err := dbContainer.MappedPort(ctx, "5432/tcp")
		if err != nil {
			_ = dbContainer.Terminate(ctx)
			setupErr = fmt.Errorf("failed to get container port: %w", err)
			return
		}

		_ = os.Setenv("DB_HOST", host)
		_ = os.Setenv("DB_PORT", port.Port())
		_ = os.Setenv("DB_DATABASE", dbName)
		_ = os.Setenv("DB_USERNAME", dbUser)
		_ = os.Setenv("DB_PASSWORD", dbPwd)
		_ = os.Setenv("DB_SCHEMA", "public")
		_ = os.Setenv("PORT", "8080")

		dbHost = host
		dbPort = port.Port()
		dbDatabase = dbName
		dbUsername = dbUser
		dbPassword = dbPwd
		dbSchema = "public"

		// Reset database instance to ensure it uses new env variables
		database.ResetInstance()

		dbService := database.New()
		testDB = dbService.DB()

		if err := runMigrationFile(ctx, testDB, initMigrationFile); err != nil {
			_ = dbContainer.Terminate(ctx)
			setupErr = fmt.Errorf("failed to run migrations: %w", err)
			return
		}

		cleanupFunc = func() error {
			_ = dbService.Close()
			return dbContainer.Terminate(ctx)
		}
	})

	return testDB, cleanupFunc, setupErr
}

// ResetTestData truncates tables and re-seeds test data.
func ResetTestData(t *testing.T) {
	t.Helper()

	db := RequireTestDB(t)
	ctx := context.Background()

	_, err := db.ExecContext(ctx, "TRUNCATE ledger, transactions, accounts, users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	if err := runMigrationFile(ctx, db, testDataMigrationFile); err != nil {
		t.Fatalf("seed test data: %v", err)
	}
}

// RequireTestDB ensures the test database is ready.
func RequireTestDB(t *testing.T) *bun.DB {
	t.Helper()
	db, _, err := SetupTestDB()
	if err != nil {
		t.Fatalf("setup test db: %v", err)
	}
	return db
}

// SetupTestServer returns an http.Handler with registered routes.
func SetupTestServer(t *testing.T) http.Handler {
	t.Helper()
	_ = RequireTestDB(t)
	return server.NewServer().Handler
}

// GetUserIDByEmail returns a user's ID by email.
func GetUserIDByEmail(t *testing.T, db *bun.DB, email string) uuid.UUID {
	t.Helper()
	var row struct {
		ID uuid.UUID `bun:"id"`
	}
	err := db.NewSelect().Table("users").Column("id").Where("email = ?", email).Scan(context.Background(), &row)
	if err != nil {
		t.Fatalf("get user id: %v", err)
	}
	return row.ID
}

// GetAccountID returns an account ID by user ID and currency.
func GetAccountID(t *testing.T, db *bun.DB, userID uuid.UUID, currency string) uuid.UUID {
	t.Helper()
	var row struct {
		ID uuid.UUID `bun:"id"`
	}
	err := db.NewSelect().Table("accounts").Column("id").Where("user_id = ?", userID).Where("currency = ?", currency).Scan(context.Background(), &row)
	if err != nil {
		t.Fatalf("get account id: %v", err)
	}
	return row.ID
}

func runMigrationFile(ctx context.Context, db *bun.DB, fileName string) error {
	root, err := repoRoot()
	if err != nil {
		return err
	}

	path := filepath.Join(root, "migrations", fileName)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sqlText := extractGooseUp(string(content))
	statements := splitSQLStatements(sqlText)
	for _, stmt := range statements {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("exec %s: %w", fileName, err)
		}
	}
	return nil
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func extractGooseUp(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	inUp := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "-- +goose Up") {
			inUp = true
			continue
		}
		if strings.HasPrefix(line, "-- +goose Down") {
			break
		}
		if inUp {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func splitSQLStatements(sqlText string) []string {
	// Split by semicolon, but preserve multi-line statements
	lines := strings.Split(sqlText, "\n")
	var statements []string
	var currentStmt strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and comment-only lines
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		// Check if line ends with semicolon (statement complete)
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			stmt := strings.TrimSpace(currentStmt.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStmt.Reset()
		}
	}

	// Handle any remaining statement without trailing semicolon
	if currentStmt.Len() > 0 {
		stmt := strings.TrimSpace(currentStmt.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}
