package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/WiggiLi/file-sharing-api/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/WiggiLi/file-sharing-api/config"
	libError "github.com/WiggiLi/file-sharing-api/lib/error"
	"github.com/WiggiLi/file-sharing-api/lib/validator"
	"github.com/WiggiLi/file-sharing-api/lib/logger"
)

type Controller struct {
	User   UserController	
	File   FileController
}

func NewController(ctx context.Context, serviceManager *service.Manager, l *logger.Logger) Controller {
	return Controller{
		User:    NewUsers(ctx, serviceManager, l),
		File:    NewFiles(ctx, serviceManager, l),
	}
}

// Start initializes Web Server, starts application and begins serving
func (controller *Controller) Start(errc chan<- error) {
	cfg := config.Get()

	// Initialize Echo instance
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.HTTPErrorHandler = libError.Error

	// Disable Echo JSON logger in debug mode
	if cfg.LogLevel == "debug" {
		if l, ok := e.Logger.(*echoLog.Logger); ok {
			l.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line}")
		}
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")
	//v1.Use(middleware.JWT([]byte("secret")))		//!DELETE
	v1.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// Auth routes
	authGroup := v1.Group("/auth")
	authGroup.POST("/register", controller.User.Register)
	authGroup.POST("/login", controller.User.Login)
	authGroup.POST("/logout", controller.User.Logout)

	authGroup.GET("/user/:id", controller.User.GetFileNamesByUserID)

	// File routes
	fileRoutes := v1.Group("/file")

	fileRoutes.POST("/", controller.File.Create) //post filename, return ID
	fileRoutes.PUT("/:id", controller.File.Upload)	//upload file
	fileRoutes.GET("/:id/meta", controller.File.Get)
	fileRoutes.DELETE("/:id/meta", controller.File.Delete)
	fileRoutes.GET("/:id", controller.File.Download)

	// Start server
	s := &http.Server{
		Addr:         cfg.HTTPAddr,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}