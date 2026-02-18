package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Password      string         `json:"-" gorm:"not null"`
	Name          string         `json:"name"`
	Avatar        string         `json:"avatar"`
	Roles         []string       `json:"roles" gorm:"type:text[]"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Session struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       string         `json:"user_id" gorm:"type:uuid;not null;index"`
	Token        string         `json:"token" gorm:"uniqueIndex;not null"`
	RefreshToken string         `json:"refresh_token" gorm:"uniqueIndex"`
	IPAddress    string         `json:"ip_address"`
	UserAgent    string         `json:"user_agent"`
	ExpiresAt    time.Time      `json:"expires_at"`
	IsRevoked    bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type Content struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Body        string         `json:"body" gorm:"type:text"`
	Excerpt     string         `json:"excerpt"`
	Status      string         `json:"status" gorm:"default:'draft'"` // draft, published, archived
	ContentType string         `json:"content_type" gorm:"default:'page'"`
	AuthorID    string         `json:"author_id" gorm:"type:uuid"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Content) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type Integration struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`     // oauth, api_key, webhook
	Provider    string         `json:"provider" gorm:"not null"` // google, github, stripe, etc.
	Credentials string         `json:"-" gorm:"type:text"`       // encrypted
	Status      string         `json:"status" gorm:"default:'active'"`
	UserID      string         `json:"user_id" gorm:"type:uuid;index"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (i *Integration) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}

type Workflow struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	TriggerType string         `json:"trigger_type" gorm:"not null"` // webhook, schedule, event
	ActionType  string         `json:"action_type" gorm:"not null"`  // http, email, notification
	Config      string         `json:"config" gorm:"type:jsonb"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	Schedule    string         `json:"schedule"` // cron expression
	LastRunAt   *time.Time     `json:"last_run_at"`
	NextRunAt   *time.Time     `json:"next_run_at"`
	UserID      string         `json:"user_id" gorm:"type:uuid;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}
