package database

import (
	"context"
	"time"

	"biometrics/pkg/models"

	"biometrics/pkg/utils"
	"github.com/google/uuid"
)

type SeedManager struct {
	db     *Postgres
	logger utils.Logger
}

func NewSeedManager(db *Postgres, logger utils.Logger) *SeedManager {
	return &SeedManager{
		db:     db,
		logger: logger,
	}
}

func (s *SeedManager) RunSeed(ctx context.Context) error {
	s.logger.Info("Starting database seeding...")

	if err := s.seedUsers(ctx); err != nil {
		return err
	}

	s.logger.Info("Database seeding completed successfully")
	return nil
}

func (s *SeedManager) seedUsers(ctx context.Context) error {
	s.logger.Info("Seeding users...")

	users := []struct {
		email         string
		firstName     string
		lastName      string
		role          string
		status        string
		emailVerified bool
	}{
		{"admin@biometrics.dev", "Admin", "User", "admin", "active", true},
		{"operator@biometrics.dev", "Operator", "User", "operator", "active", true},
		{"auditor@biometrics.dev", "Auditor", "User", "auditor", "active", true},
		{"user@biometrics.dev", "Test", "User", "user", "active", true},
		{"guest@biometrics.dev", "Guest", "User", "guest", "active", false},
	}

	for _, u := range users {
		var existing models.User
		err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", u.email).Scan(&existing.ID)
		if err == nil {
			s.logger.Info("User already exists, skipping", "email", u.email)
			continue
		}

		now := time.Now()
		passwordHash, err := utils.HashPassword("password123")
		if err != nil {
			s.logger.Error("Failed to hash password", "error", err)
			continue
		}

		user := &models.User{
			ID:                  uuid.New(),
			Email:               u.email,
			PasswordHash:        passwordHash,
			FirstName:           u.firstName,
			LastName:            u.lastName,
			Phone:               "",
			Role:                u.role,
			Status:              u.status,
			EmailVerified:       u.emailVerified,
			TwoFactorEnabled:    false,
			Provider:            "local",
			Avatar:              "",
			Bio:                 "",
			Roles:               []string{u.role},
			Permissions:         []string{"read"},
			FailedLoginAttempts: 0,
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		if err := s.db.CreateUser(ctx, user); err != nil {
			s.logger.Error("Failed to create user", "email", u.email, "error", err)
			continue
		}

		s.logger.Info("User created", "email", u.email, "role", u.role)
	}

	return nil
}

func (s *SeedManager) seedDemoData(ctx context.Context) error {
	s.logger.Info("Seeding demo content...")

	var adminID uuid.UUID
	err := s.db.QueryRow(ctx, "SELECT id FROM users WHERE role = 'admin' LIMIT 1").Scan(&adminID)
	if err != nil {
		return err
	}

	demoContent := []struct {
		title       string
		slug        string
		body        string
		contentType string
		status      string
	}{
		{
			title:       "Welcome to BIOMETRICS",
			slug:        "welcome-to-biometrics",
			body:        "# Welcome\n\nThis is a demo content for the BIOMETRICS platform.",
			contentType: "page",
			status:      "published",
		},
		{
			title:       "Getting Started Guide",
			slug:        "getting-started-guide",
			body:        "# Getting Started\n\nLearn how to use the BIOMETRICS platform.",
			contentType: "article",
			status:      "published",
		},
		{
			title:       "API Documentation",
			slug:        "api-documentation",
			body:        "# API Documentation\n\nComplete API reference for BIOMETRICS.",
			contentType: "article",
			status:      "draft",
		},
	}

	for _, c := range demoContent {
		var existingID uuid.UUID
		err := s.db.QueryRow(ctx, "SELECT id FROM content WHERE slug = $1", c.slug).Scan(&existingID)
		if err == nil {
			s.logger.Info("Content already exists, skipping", "slug", c.slug)
			continue
		}

		now := time.Now()
		_, err = s.db.Exec(ctx, `INSERT INTO content 
			(id, tenant_id, title, slug, body, type, status, author_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			uuid.New(), nil, c.title, c.slug, c.body, c.contentType, c.status, adminID, now, now)

		if err != nil {
			s.logger.Error("Failed to create content", "slug", c.slug, "error", err)
			continue
		}

		s.logger.Info("Content created", "slug", c.slug)
	}

	return nil
}
