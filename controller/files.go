package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	b64 "encoding/base64"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/WiggiLi/file-sharing-api/lib/logger"
	"github.com/WiggiLi/file-sharing-api/lib/types"
	"github.com/WiggiLi/file-sharing-api/service"
	"github.com/WiggiLi/file-sharing-api/model"
)

// FileController ...
type FileController struct {
	ctx      context.Context
	services *service.Manager
	logger   *logger.Logger
}

// NewFiles creates a new file controller.
func NewFiles(ctx context.Context, services *service.Manager, logger *logger.Logger) FileController {
	return FileController{
		ctx:      ctx,
		services: services,
		logger:   logger,
	}
}

// Create creates new file
func (ctr *FileController) Create(ctx echo.Context) error {
	var file model.File
	err := ctx.Bind(&file)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode file data"))
	}
	err = ctx.Validate(&file)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	//check if user login
	sess, _ := session.Get("session", ctx)
	var UserID interface{} = sess.Values["id_user"]
	var UserUUID uuid.UUID
	if UserID != nil { 
		uUUID := sess.Values["id_user"].(string) 
		UserUUID, err = uuid.Parse(uUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse file UUID"))
		}
	} 

	createdFile, err := ctr.services.FileMeta.CreateFileMeta(ctx.Request().Context(), &file, UserUUID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not create file"))
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create file"))
		}
	}

	ctr.logger.Debug().Msgf("Created file '%s'", createdFile.ID.String())
	return ctx.JSON(http.StatusCreated, createdFile)
}

// returns file by ID
func (ctr *FileController) Get(ctx echo.Context) error {
	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse file UUID"))
	}
	file, err := ctr.services.FileMeta.GetFileMeta(ctx.Request().Context(), fileID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get file"))
		}
	}
	return ctx.JSON(http.StatusOK, file)
}

// Delete deletes file by ID
func (ctr *FileController) Delete(ctx echo.Context) error {
	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse file UUID"))
	}
	err = ctr.services.FileMeta.DeleteFileMeta(ctx.Request().Context(), fileID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, errors.Wrap(err, "could not delete file"))
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not delete file"))
		}
	}

	ctr.logger.Debug().Msgf("Deleted file '%s'", fileID.String())

	return ctx.JSON(http.StatusOK, map[string]bool{"ok": true,})
}

// Upload file content
func (ctr *FileController) Upload(ctx echo.Context) error {
	fileBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not read file body"))
	}
	defer ctx.Request().Body.Close()

	fileID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse file UUID"))
	}
	err = ctr.services.FileContent.Upload(ctx.Request().Context(), fileID, fileBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not upload file"))
	}
	ctr.logger.Debug().Msgf("Saved content for file '%s'", fileID.String())
	return ctx.JSON(http.StatusOK, map[string]bool{"ok": true,})
}

// Download file content from the cloud
func (ctr *FileController) Download(ctx echo.Context) error {
	link := ctx.Param("id")
	sDec, _ := b64.StdEncoding.DecodeString(link)
	fileID, err := uuid.Parse(string(sDec))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse file UUID"))
	}
	fileBody, dbFile, err := ctr.services.FileContent.Download(ctx.Request().Context(), fileID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not download file content"))
	}
	ctr.logger.Debug().Msgf("Downloaded content for file '%s'", fileID.String())
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", dbFile.Filename))
	return ctx.Blob(http.StatusOK, "", fileBody)
}
