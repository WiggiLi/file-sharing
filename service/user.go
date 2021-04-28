package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/WiggiLi/file-sharing-api/lib/types"
	"github.com/WiggiLi/file-sharing-api/model"
	"github.com/WiggiLi/file-sharing-api/store"
)

type UserWebService struct {
	ctx context.Context
	store *store.Store
}

func NewUserWebService(ctx context.Context, store *store.Store) *UserWebService {
	return &UserWebService{
		ctx: ctx,
		store: store,
	}
}

func (svc *UserWebService) GetFileNamesByUserID(ctx context.Context, userID uuid.UUID) (*[]model.File, error) {
	/*
	//get from cache, if possible
	userCache, err := svc.store.User.Redis.GetFileNamesByUserID(ctx, fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, errors.Wrap(err, "newsUC.GetNewsByID.GetNewsByIDCtx.")
	}
	if userCache != nil {
		return userCache, nil
	}
	*/

	//get from DB
	userDB, err := svc.store.User.Pg.GetFileNamesByUserID(ctx, userID)

	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser")
	}
	if userDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%s' not found", userID.String()))
	}

	//add to cache
	err = svc.store.User.Redis.CreateFileNamesByUserID(ctx, fmt.Sprintf("%d", userID), 3600, userDB)
	if err != nil {
		return nil, errors.Wrap(err, "newsUC.GetNewsByID.GetNewsByIDCtx.")
	}

	return userDB, nil 	
}

// CreateUser ...
func (svc UserWebService) CreateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {
	reqUser.ID = uuid.New()

	_, err := svc.store.User.Pg.CreateUser(ctx, reqUser)	 	
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.CreateUser error")
	}

	return reqUser, nil
}

// CreateUser ...
func (svc UserWebService) LoginUser(ctx context.Context, reqUser *model.User) (*model.User, error) {
	reqUser.ID = uuid.New()

	newUser, err := svc.store.User.Pg.LoginUser(ctx, reqUser)	 	
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.LoginUser error")
	}

	return newUser, nil
}
