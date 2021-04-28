package service

import (
	"context"
	"fmt"
	b64 "encoding/base64"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/WiggiLi/file-sharing-api/lib/types"
	"github.com/WiggiLi/file-sharing-api/model"
	"github.com/WiggiLi/file-sharing-api/store"
)

type FileMetaSvc struct {
	ctx   context.Context
	store *store.Store
}

// NewFileMetaSvc creates a new file web service
func NewFileMetaSvc(ctx context.Context, store *store.Store) *FileMetaSvc {
	return &FileMetaSvc{
		ctx:   ctx,
		store: store,
	}
}

func (svc *FileMetaSvc) GetFileMeta(ctx context.Context, fileID uuid.UUID) (*model.File, error) {
	fileDB, err := svc.store.File.GetFileMeta(ctx, fileID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.GetFileMeta")
	}
	if fileDB == nil {
		return nil, errors.Wrap(types.ErrBadRequest, fmt.Sprintf("File '%s' not found", fileID.String()))
	}

	return fileDB, nil
}

func (svc FileMetaSvc) CreateFileMeta(ctx context.Context, reqFile *model.File, authenticated uuid.UUID) (*model.File, error) {
	reqFile.ID = uuid.New()
	reqFile.Link = b64.StdEncoding.EncodeToString([]byte(reqFile.ID.String()))

	_, err := svc.store.File.CreateFileMeta(ctx, reqFile, authenticated)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.CreateFileMeta error")
	}

	// get created file by ID
	createdDBFile, err := svc.store.File.GetFileMeta(ctx, reqFile.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.GetFileMeta error")
	}

	return createdDBFile, nil		//?  return file back?
}

func (svc *FileMetaSvc) UpdateFileMeta(ctx context.Context, reqFile *model.File) (*model.File, error) {
	fileDB, err := svc.store.File.GetFileMeta(ctx, reqFile.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.GetFileMeta error")
	}
	if fileDB == nil {
		return nil, errors.Wrap(types.ErrBadRequest, fmt.Sprintf("File '%s' not found", reqFile.ID.String()))
	}

	// update file
	_, err = svc.store.File.UpdateFileMeta(ctx, reqFile)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.UpdateFileMeta error")
	}

	// get updated file by ID
	updatedDBFile, err := svc.store.File.GetFileMeta(ctx, reqFile.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.file.GetFileMeta error")
	}

	return updatedDBFile, nil
}

func (svc *FileMetaSvc) DeleteFileMeta(ctx context.Context, fileID uuid.UUID) error {
	// Check if file exists
	fileDB, err := svc.store.File.GetFileMeta(ctx, fileID)
	if err != nil {
		return errors.Wrap(err, "svc.file.GetFileMeta error")
	}
	if fileDB == nil {
		return errors.Wrap(types.ErrNotFound, fmt.Sprintf("File '%s' not found", fileID.String()))
	}

	err = svc.store.File.DeleteFileMeta(ctx, fileID)
	if err != nil {
		return errors.Wrap(err, "svc.file.DeleteFileMeta error")
	}

	return nil
}
