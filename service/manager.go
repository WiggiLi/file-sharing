package service

import (
	"context"

	"github.com/WiggiLi/file-sharing-api/store"
	"github.com/pkg/errors"
	"github.com/WiggiLi/file-sharing-api/model"
)

type Manager struct {
	User model.UserService
	FileMeta model.FileMetaService
	FileContent model.FileContentService
}

func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}
	return &Manager{
		User:        NewUserWebService(ctx, store),
		FileMeta:    NewFileMetaSvc(ctx, store),
		FileContent: NewFileContentSvc(ctx, store),
	}, nil
}