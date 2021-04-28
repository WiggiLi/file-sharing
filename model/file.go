package model

import (
	"time"
	"context"
	"github.com/google/uuid"
)

// File holds file metadata as a JSON
type File struct {
	ID        uuid.UUID `json:"id" pg:"id,notnull,pk"`
	Filename  string    `json:"filename" validate:"required" pg:"filename,notnull"`
	Link	  string	`json:"link" pg:"link,notnull"`
	CreatedAt time.Time `json:"created_at" pg:"created_at,notnull"`
}


// FileMetaService is a service for files
//go:generate mockery --dir . --name FileMetaService --output ./mocks
type FileMetaService interface {
	GetFileMeta(context.Context, uuid.UUID) (*File, error)
	CreateFileMeta(context.Context, *File, uuid.UUID) (*File, error)
	UpdateFileMeta(context.Context, *File) (*File, error)
	DeleteFileMeta(context.Context, uuid.UUID) error
}

// FileContentService is a service to upload file content
//go:generate mockery --dir . --name FileContentService --output ./mocks
type FileContentService interface {
	Upload(context.Context, uuid.UUID, []byte) error
	Download(context.Context, uuid.UUID) ([]byte, *File, error)
}

// FileMetaRepo is a store for files
//go:generate mockery --dir . --name FileMetaRepo --output ./mocks
type FileMetaRepo interface {
	GetFileMeta(context.Context, uuid.UUID) (*File, error)
	CreateFileMeta(context.Context, *File, uuid.UUID) (*File, error)
	UpdateFileMeta(context.Context, *File) (*File, error)
	DeleteFileMeta(context.Context, uuid.UUID) error
}

// FileContentRepo is a store for file content
//go:generate mockery --dir . --name FileContentRepo --output ./mocks
type FileContentRepo interface {
	Upload(context.Context, *File, []byte) error
	Download(context.Context, *File) ([]byte, error)
}