package models

import (
	"time"

	"github.com/google/uuid"
)

type Content struct {
	ID              uuid.UUID              `json:"id" db:"id"`
	TenantID        *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Title           string                 `json:"title" db:"title"`
	Slug            string                 `json:"slug" db:"slug"`
	Body            string                 `json:"body" db:"body"`
	Excerpt         string                 `json:"excerpt" db:"excerpt"`
	Type            string                 `json:"type" db:"type"`
	Status          string                 `json:"status" db:"status"`
	AuthorID        uuid.UUID              `json:"author_id" db:"author_id"`
	FeaturedImage   string                 `json:"featured_image" db:"featured_image"`
	MetaTitle       string                 `json:"meta_title" db:"meta_title"`
	MetaDescription string                 `json:"meta_description" db:"meta_description"`
	Tags            []string               `json:"tags" db:"tags"`
	Categories      []string               `json:"categories" db:"categories"`
	PublishedAt     *time.Time             `json:"published_at" db:"published_at"`
	ScheduledAt     *time.Time             `json:"scheduled_at" db:"scheduled_at"`
	ViewCount       int                    `json:"view_count" db:"view_count"`
	LikeCount       int                    `json:"like_count" db:"like_count"`
	ShareCount      int                    `json:"share_count" db:"share_count"`
	IsFeatured      bool                   `json:"is_featured" db:"is_featured"`
	IsPublic        bool                   `json:"is_public" db:"is_public"`
	Language        string                 `json:"language" db:"language"`
	Version         int                    `json:"version" db:"version"`
	ParentID        *uuid.UUID             `json:"parent_id,omitempty" db:"parent_id"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type ContentType string

const (
	ContentTypeArticle  ContentType = "article"
	ContentTypePage     ContentType = "page"
	ContentTypePost     ContentType = "post"
	ContentTypeDocument ContentType = "document"
	ContentTypeMedia    ContentType = "media"
)

type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusScheduled ContentStatus = "scheduled"
	ContentStatusArchived  ContentStatus = "archived"
	ContentStatusReview    ContentStatus = "review"
)

type CreateContentInput struct {
	TenantID        string   `json:"tenant_id"`
	Title           string   `json:"title" binding:"required"`
	Slug            string   `json:"slug"`
	Body            string   `json:"body" binding:"required"`
	Excerpt         string   `json:"excerpt"`
	Type            string   `json:"type" binding:"required"`
	Status          string   `json:"status"`
	AuthorID        string   `json:"author_id" binding:"required"`
	FeaturedImage   string   `json:"featured_image"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Tags            []string `json:"tags"`
	Categories      []string `json:"categories"`
	PublishedAt     string   `json:"published_at"`
	Language        string   `json:"language"`
	ParentID        string   `json:"parent_id"`
}

type UpdateContentInput struct {
	Title           *string  `json:"title,omitempty"`
	Slug            *string  `json:"slug,omitempty"`
	Body            *string  `json:"body,omitempty"`
	Excerpt         *string  `json:"excerpt,omitempty"`
	Status          *string  `json:"status,omitempty"`
	FeaturedImage   *string  `json:"featured_image,omitempty"`
	MetaTitle       *string  `json:"meta_title,omitempty"`
	MetaDescription *string  `json:"meta_description,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Categories      []string `json:"categories,omitempty"`
	PublishedAt     *string  `json:"published_at,omitempty"`
	IsFeatured      *bool    `json:"is_featured,omitempty"`
	IsPublic        *bool    `json:"is_public,omitempty"`
}

type ContentFilter struct {
	Type     string `query:"type"`
	Status   string `query:"status"`
	AuthorID string `query:"author_id"`
	TenantID string `query:"tenant_id"`
	Tag      string `query:"tag"`
	Category string `query:"category"`
	Search   string `query:"search"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

func (c *Content) ToResponse() *ContentResponse {
	return &ContentResponse{
		ID:            c.ID,
		TenantID:      c.TenantID,
		Title:         c.Title,
		Slug:          c.Slug,
		Body:          c.Body,
		Excerpt:       c.Excerpt,
		Type:          c.Type,
		Status:        c.Status,
		AuthorID:      c.AuthorID,
		FeaturedImage: c.FeaturedImage,
		Tags:          c.Tags,
		Categories:    c.Categories,
		PublishedAt:   c.PublishedAt,
		ViewCount:     c.ViewCount,
		LikeCount:     c.LikeCount,
		ShareCount:    c.ShareCount,
		IsFeatured:    c.IsFeatured,
		Language:      c.Language,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

type ContentResponse struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Body          string     `json:"body"`
	Excerpt       string     `json:"excerpt"`
	Type          string     `json:"type"`
	Status        string     `json:"status"`
	AuthorID      uuid.UUID  `json:"author_id"`
	FeaturedImage string     `json:"featured_image"`
	Tags          []string   `json:"tags"`
	Categories    []string   `json:"categories"`
	PublishedAt   *time.Time `json:"published_at"`
	ViewCount     int        `json:"view_count"`
	LikeCount     int        `json:"like_count"`
	ShareCount    int        `json:"share_count"`
	IsFeatured    bool       `json:"is_featured"`
	Language      string     `json:"language"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
