package controller

import (
	"context"
	"net/http"
	"log"		//! DELTE
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/WiggiLi/file-sharing-api/lib/logger"
	"github.com/WiggiLi/file-sharing-api/lib/types"
	"github.com/WiggiLi/file-sharing-api/service"
	"github.com/WiggiLi/file-sharing-api/model"
)

// UserController ...
type UserController struct {
	ctx      context.Context
	services *service.Manager
	logger   *logger.Logger
}

// NewUsers creates a new user controller.
func NewUsers(ctx context.Context, services *service.Manager, logger *logger.Logger) UserController {
	return UserController{
		ctx:      ctx,
		services: services,
		logger:   logger,
	}
}

// Create creates new user
func (ctr *UserController) Register(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	//ctx.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var user model.User
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode user data"))
	}
	err = ctx.Validate(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	createdUser, err := ctr.services.User.CreateUser(ctx.Request().Context(), &user)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		case errors.Cause(err) == types.ErrUnauthorized:
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create user"))
		}
	}

	//session
	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	  }
	sess.Values["id_user"] = createdUser.ID.String()
	err = sess.Save(ctx.Request(), ctx.Response())
	if err != nil {
		ctr.logger.Debug().Msgf("Session save error, ", err)
	}

	ctr.logger.Debug().Msgf("Created user '%s'", createdUser.ID.String())

	return ctx.JSON(http.StatusCreated, createdUser) //createdUser
}

// Login user
func (ctr *UserController) Login(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	//ctx.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var user model.User
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode user data"))
	}
	err = ctx.Validate(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}


	createdUser, err := ctr.services.User.LoginUser(ctx.Request().Context(), &user)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		case errors.Cause(err) == types.ErrUnauthorized:
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create user"))
		}
	}

	sess, _ := session.Get("session", ctx)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	  }
	//sess.Values["authenticated"] = true
	//gob.Register(uuid.UUID)
	sess.Values["id_user"] = createdUser.ID.String()
	err = sess.Save(ctx.Request(), ctx.Response())
	if err != nil {
		ctr.logger.Debug().Msgf("Session save error, ", err)
	}

	ctr.logger.Debug().Msgf("Created user '%s'", createdUser.ID.String())

	return ctx.JSON(http.StatusCreated, createdUser) //createdUser
}

func (ctr *UserController) Logout(ctx echo.Context) error {
	sess, _ := session.Get("session", ctx)
	sess.Values["id_user"] = uuid.Nil
	sess.Save(ctx.Request(), ctx.Response())

	return ctx.NoContent(http.StatusOK)
}

// Get returns user by ID
func (ctr *UserController) GetFileNamesByUserID(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse user UUID"))
	}
	
	//session
	sess, _ := session.Get("session", ctx)

	var UserID uuid.UUID = sess.Values["id_user"].(uuid.UUID) 
	if UserID == uuid.Nil { 
		log.Println("authenticated Nil")
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	} 


	fileNames, err := ctr.services.User.GetFileNamesByUserID(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			log.Println("UserController2")
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Cause(err) == types.ErrUnauthorized:
			return echo.NewHTTPError(http.StatusUnauthorized, err)	
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get user"))
		}
	}
	
	return ctx.JSON(http.StatusOK, fileNames) 
}

