package database

import (
	"context"
	"fmt"
	"time"

	"biometrics/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DB struct {
	*Postgres
}

func NewDB(postgres *Postgres) *DB {
	return &DB{Postgres: postgres}
}

func (d *DB) GormDB() *gorm.DB {
	return nil
}

type UserRepository struct {
	db *Postgres
}

func NewUserRepository(db *Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.CreateUser(ctx, user)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return r.db.GetUserByID(ctx, id.String())
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.db.GetUserByEmail(ctx, email)
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.UpdateUser(ctx, user)
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteUser(ctx, id.String())
}

func (r *UserRepository) List(ctx context.Context, filter models.UserFilter) ([]*models.User, int64, error) {
	return r.db.ListUsers(ctx, filter)
}

type ContentRepository struct {
	db *Postgres
}

func NewContentRepository(db *Postgres) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(ctx context.Context, content *models.Content) error {
	sql := `INSERT INTO content (id, tenant_id, title, slug, body, excerpt, type, status, author_id,
			featured_image, meta_title, meta_description, tags, categories, published_at, scheduled_at,
			view_count, like_count, share_count, is_featured, is_public, language, version, parent_id,
			metadata, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)`

	_, err := r.db.Exec(ctx, sql,
		content.ID, content.TenantID, content.Title, content.Slug, content.Body, content.Excerpt,
		content.Type, content.Status, content.AuthorID, content.FeaturedImage, content.MetaTitle,
		content.MetaDescription, content.Tags, content.Categories, content.PublishedAt, content.ScheduledAt,
		content.ViewCount, content.LikeCount, content.ShareCount, content.IsFeatured, content.IsPublic,
		content.Language, content.Version, content.ParentID, content.Metadata, content.CreatedAt, content.UpdatedAt)

	return err
}

func (r *ContentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Content, error) {
	sql := `SELECT id, tenant_id, title, slug, body, excerpt, type, status, author_id,
			featured_image, meta_title, meta_description, tags, categories, published_at, scheduled_at,
			view_count, like_count, share_count, is_featured, is_public, language, version, parent_id,
			metadata, created_at, updated_at, deleted_at
			FROM content WHERE id = $1 AND deleted_at IS NULL`

	row := r.db.QueryRow(ctx, sql, id)

	var content models.Content
	err := row.Scan(
		&content.ID, &content.TenantID, &content.Title, &content.Slug, &content.Body, &content.Excerpt,
		&content.Type, &content.Status, &content.AuthorID, &content.FeaturedImage, &content.MetaTitle,
		&content.MetaDescription, &content.Tags, &content.Categories, &content.PublishedAt, &content.ScheduledAt,
		&content.ViewCount, &content.LikeCount, &content.ShareCount, &content.IsFeatured, &content.IsPublic,
		&content.Language, &content.Version, &content.ParentID, &content.Metadata, &content.CreatedAt,
		&content.UpdatedAt, &content.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (r *ContentRepository) Update(ctx context.Context, content *models.Content) error {
	sql := `UPDATE content SET 
			title = $1, slug = $2, body = $3, excerpt = $4, status = $5,
			featured_image = $6, meta_title = $7, meta_description = $8,
			tags = $9, categories = $10, published_at = $11, is_featured = $12,
			is_public = $13, language = $14, version = version + 1, updated_at = $15
			WHERE id = $16`

	_, err := r.db.Exec(ctx, sql,
		content.Title, content.Slug, content.Body, content.Excerpt, content.Status,
		content.FeaturedImage, content.MetaTitle, content.MetaDescription,
		content.Tags, content.Categories, content.PublishedAt, content.IsFeatured,
		content.IsPublic, content.Language, content.UpdatedAt, content.ID)

	return err
}

func (r *ContentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE content SET deleted_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, sql, time.Now(), id)
	return err
}

func (r *ContentRepository) List(ctx context.Context, filter models.ContentFilter) ([]*models.Content, int64, error) {
	baseSQL := `FROM content WHERE deleted_at IS NULL`
	countSQL := `SELECT COUNT(*) ` + baseSQL
	selectSQL := `SELECT id, tenant_id, title, slug, body, excerpt, type, status, author_id,
			featured_image, meta_title, meta_description, tags, categories, published_at, scheduled_at,
			view_count, like_count, share_count, is_featured, is_public, language, version, parent_id,
			metadata, created_at, updated_at, deleted_at ` + baseSQL

	var args []interface{}
	argIndex := 1

	if filter.Type != "" {
		selectSQL += fmt.Sprintf(" AND type = $%d", argIndex)
		countSQL += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Status != "" {
		selectSQL += fmt.Sprintf(" AND status = $%d", argIndex)
		countSQL += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	var total int64
	r.db.QueryRow(ctx, countSQL, args...).Scan(&total)

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize
	selectSQL += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(
			&content.ID, &content.TenantID, &content.Title, &content.Slug, &content.Body, &content.Excerpt,
			&content.Type, &content.Status, &content.AuthorID, &content.FeaturedImage, &content.MetaTitle,
			&content.MetaDescription, &content.Tags, &content.Categories, &content.PublishedAt, &content.ScheduledAt,
			&content.ViewCount, &content.LikeCount, &content.ShareCount, &content.IsFeatured, &content.IsPublic,
			&content.Language, &content.Version, &content.ParentID, &content.Metadata, &content.CreatedAt,
			&content.UpdatedAt, &content.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		contents = append(contents, &content)
	}

	return contents, total, nil
}

type BiometricRepository struct {
	db *Postgres
}

func NewBiometricRepository(db *Postgres) *BiometricRepository {
	return &BiometricRepository{db: db}
}

func (r *BiometricRepository) Create(ctx context.Context, biometric *models.Biometric) error {
	sql := `INSERT INTO biometrics (id, user_id, tenant_id, type, sub_type, label, template_encrypted,
			algorithm, device_id, device_info, public_key, credential_id, status, enrolled_at,
			last_used_at, expires_at, failure_count, max_failures, verify_count, success_count,
			quality_score, metadata, is_active, is_primary, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`

	_, err := r.db.Exec(ctx, sql,
		biometric.ID, biometric.UserID, biometric.TenantID, biometric.Type, biometric.SubType,
		biometric.Label, biometric.Template, biometric.Algorithm, biometric.DeviceID,
		biometric.DeviceInfo, biometric.PublicKey, biometric.CredentialID, biometric.Status,
		biometric.EnrolledAt, biometric.LastUsedAt, biometric.ExpiresAt, biometric.FailureCount,
		biometric.MaxFailures, biometric.VerifyCount, biometric.SuccessCount, biometric.QualityScore,
		biometric.Metadata, biometric.IsActive, biometric.IsPrimary, biometric.Version,
		biometric.CreatedAt, biometric.UpdatedAt)

	return err
}

func (r *BiometricRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Biometric, error) {
	sql := `SELECT id, user_id, tenant_id, type, sub_type, label, template_encrypted,
			algorithm, device_id, device_info, public_key, credential_id, status, enrolled_at,
			last_used_at, expires_at, failure_count, max_failures, verify_count, success_count,
			quality_score, metadata, is_active, is_primary, version, created_at, updated_at, deleted_at
			FROM biometrics WHERE id = $1 AND deleted_at IS NULL`

	row := r.db.QueryRow(ctx, sql, id)

	var biometric models.Biometric
	err := row.Scan(
		&biometric.ID, &biometric.UserID, &biometric.TenantID, &biometric.Type, &biometric.SubType,
		&biometric.Label, &biometric.Template, &biometric.Algorithm, &biometric.DeviceID,
		&biometric.DeviceInfo, &biometric.PublicKey, &biometric.CredentialID, &biometric.Status,
		&biometric.EnrolledAt, &biometric.LastUsedAt, &biometric.ExpiresAt, &biometric.FailureCount,
		&biometric.MaxFailures, &biometric.VerifyCount, &biometric.SuccessCount, &biometric.QualityScore,
		&biometric.Metadata, &biometric.IsActive, &biometric.IsPrimary, &biometric.Version,
		&biometric.CreatedAt, &biometric.UpdatedAt, &biometric.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &biometric, nil
}

func (r *BiometricRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Biometric, error) {
	sql := `SELECT id, user_id, tenant_id, type, sub_type, label, template_encrypted,
			algorithm, device_id, device_info, public_key, credential_id, status, enrolled_at,
			last_used_at, expires_at, failure_count, max_failures, verify_count, success_count,
			quality_score, metadata, is_active, is_primary, version, created_at, updated_at, deleted_at
			FROM biometrics WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var biometrics []*models.Biometric
	for rows.Next() {
		var biometric models.Biometric
		err := rows.Scan(
			&biometric.ID, &biometric.UserID, &biometric.TenantID, &biometric.Type, &biometric.SubType,
			&biometric.Label, &biometric.Template, &biometric.Algorithm, &biometric.DeviceID,
			&biometric.DeviceInfo, &biometric.PublicKey, &biometric.CredentialID, &biometric.Status,
			&biometric.EnrolledAt, &biometric.LastUsedAt, &biometric.ExpiresAt, &biometric.FailureCount,
			&biometric.MaxFailures, &biometric.VerifyCount, &biometric.SuccessCount, &biometric.QualityScore,
			&biometric.Metadata, &biometric.IsActive, &biometric.IsPrimary, &biometric.Version,
			&biometric.CreatedAt, &biometric.UpdatedAt, &biometric.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		biometrics = append(biometrics, &biometric)
	}

	return biometrics, nil
}

func (r *BiometricRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE biometrics SET deleted_at = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, sql, time.Now(), id)
	return err
}
