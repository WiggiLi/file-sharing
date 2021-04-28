package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/go-pg/pg/v10"
	
	"github.com/WiggiLi/file-sharing-api/model"
)

// FileMetaPgRepo ...
type FileMetaPgRepo struct {
	db *DB
}

// NewFileMetaRepo ...
func NewFileMetaRepo(db *DB) *FileMetaPgRepo {
	return &FileMetaPgRepo{db: db}
}

// GetFileMeta retrieves file from MySQL
func (repo *FileMetaPgRepo) GetFileMeta(ctx context.Context, id uuid.UUID) (*model.File, error) {
	file := &model.File{}
	err := repo.db.Model(file).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}

// CreateFileMeta creates file in Postgres
func (repo *FileMetaPgRepo) CreateFileMeta(ctx context.Context, file *model.File, authenticated uuid.UUID) (*model.File, error) {
	_, err := repo.db.Model(file).Insert()
	if err != nil {
		return nil, err
	}

	if authenticated != uuid.Nil {
		filesOfUser := &model.FilesOfUsers{} 
		filesOfUser.User_ID = authenticated
		filesOfUser.File_ID = file.ID

		_, err := repo.db.Model(filesOfUser).Insert()
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

// UpdateFileMeta updates file in Postgres
func (repo *FileMetaPgRepo) UpdateFileMeta(ctx context.Context, file *model.File) (*model.File, error) {
	_, err := repo.db.Model(file).WherePK().Update()
	if err != nil {
		if err == pg.ErrNoRows { //not found
			return nil, nil
		}
		return nil, err
	}

	return file, nil
}

// DeleteFileMeta deletes file in Postgres
func (repo *FileMetaPgRepo) DeleteFileMeta(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Model((*model.File)(nil)).Table("files").
		Where("id = ?", id).
		Delete()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
