package local

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/WiggiLi/file-sharing-api/model"
	"github.com/pkg/errors"
)


// FileContentRepo ...
type LocalFileContentRepo struct {
	filePath string
}

// NewFileContentRepo ...
func NewFileContentRepo(filePath string) *LocalFileContentRepo {
	return &LocalFileContentRepo{filePath: filePath}
}

// Upload file to Google Cloud storage
func (repo *LocalFileContentRepo) Upload(ctx context.Context, dbFile *model.File, fileBody []byte) error {
	if dbFile == nil {
		return errors.New("No DB file provided")
	}

	if len(fileBody) == 0 {
		return errors.New("No file body provided to upload")
	}

	if err := os.MkdirAll(repo.filePath, os.ModePerm); err != nil {
		return errors.Wrap(err, "os.MkdirAll failed")		//? MkdirAll?
	}

	return ioutil.WriteFile(repo.filePath+"/"+dbFile.Filename, fileBody, 0644)
}

// Download file from Google Cloud storage
func (repo *LocalFileContentRepo) Download(ctx context.Context, dbFile *model.File) ([]byte, error) {
	if dbFile == nil {
		return nil, errors.New("No DB file provided")
	}

	return ioutil.ReadFile(repo.filePath + "/" + dbFile.Filename)
}
