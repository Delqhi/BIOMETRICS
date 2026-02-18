# KNOWLEDGE-BASE.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The Knowledge Base module provides a comprehensive self-service documentation and help center for BIOMETRICS users. This document describes the architecture, features, and integration points for the knowledge base system.

## Architecture

### Components

| Component | Description | Technology |
|-----------|-------------|------------|
| Content Management | Article CRUD operations | Supabase PostgreSQL |
| Search Engine | Full-text search | Supabase Text Search |
| Categories | Hierarchical category structure | Tree structure |
| Tags | Flexible tagging system | PostgreSQL array |
| Analytics | View tracking, feedback | Analytics tables |

### Database Schema

```sql
-- Categories table
CREATE TABLE kb_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES kb_categories(id),
    icon VARCHAR(100),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Articles table
CREATE TABLE kb_articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    category_id UUID REFERENCES kb_categories(id),
    author_id UUID REFERENCES auth.users(id),
    status VARCHAR(50) DEFAULT 'draft',
    view_count INT DEFAULT 0,
    helpful_count INT DEFAULT 0,
    not_helpful_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

-- Article tags
CREATE TABLE kb_article_tags (
    article_id UUID REFERENCES kb_articles(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES kb_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

-- Tags table
CREATE TABLE kb_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Feedback table
CREATE TABLE kb_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID REFERENCES kb_articles(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id),
    helpful BOOLEAN NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## API Endpoints

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/kb/categories | List all categories |
| GET | /api/kb/categories/:slug | Get category by slug |
| POST | /api/kb/categories | Create category (admin) |
| PUT | /api/kb/categories/:id | Update category (admin) |
| DELETE | /api/kb/categories/:id | Delete category (admin) |

### Articles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/kb/articles | List articles (with pagination) |
| GET | /api/kb/articles/:slug | Get article by slug |
| GET | /api/kb/articles/search | Search articles |
| POST | /api/kb/articles | Create article (admin) |
| PUT | /api/kb/articles/:id | Update article (admin) |
| DELETE | /api/kb/articles/:id | Delete article (admin) |
| POST | /api/kb/articles/:id/feedback | Submit feedback |

### Search

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/kb/search | Full-text search |
| GET | /api/kb/popular | Get popular articles |
| GET | /api/kb/related/:id | Get related articles |

## Features

### 1. Article Management
- Rich text editor with markdown support
- Image upload and embedding
- Code snippet syntax highlighting
- Version history
- Draft/Published status

### 2. Category Organization
- Hierarchical categories (unlimited depth)
- Drag-and-drop reordering
- Category icons
- Breadcrumb navigation

### 3. Search
- Full-text search across titles and content
- Fuzzy matching
- Search suggestions
- Filter by category/tag

### 4. Analytics
- View count tracking
- User feedback collection
- Search term analytics
- Popular articles ranking

### 5. User Features
- Bookmark articles
- Print-friendly view
- Share articles
- Language support (i18n ready)

## Integration Points

### Supabase Functions
- `kb-get-categories` - Fetch category tree
- `kb-search-articles` - Full-text search
- `kb-track-view` - View count increment

### n8n Workflows
- Article publish notifications
- Feedback alert escalation
- Weekly analytics report

## Security

- Row Level Security (RLS) enabled
- Admin role required for write operations
- Public read access for published articles
- Rate limiting on search endpoints

## Performance

- CDN caching for static content
- Database indexes on frequently queried columns
- Pagination for large lists
- Optimistic UI updates

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
