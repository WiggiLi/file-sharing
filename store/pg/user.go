package pg

import (
	"context"
	"time"
	"log"
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"

	"github.com/WiggiLi/file-sharing-api/model"
)

// UserPgRepo ...
type UserPgRepo struct {
	db *DB
}

// NewUserRepo ...
func NewUserRepo(db *DB) *UserPgRepo {
	return &UserPgRepo{db: db}
}

// GetUser retrieves user from Postgres
func (repo *UserPgRepo) GetFileNamesByUserID(ctx context.Context, id uuid.UUID) (*[]model.File, error) {
	user := &[]model.File{}

	authorBooks := repo.db.Model((*model.FilesOfUsers)(nil)).ColumnExpr("id_file").Where("id_user = ?::uuid", id)

	err := repo.db.Model(user).Where("id IN ( ?::uuid )", authorBooks).Select()

	if err != nil {
		if err == pg.ErrNoRows { //not found
			log.Println("ErrNoRows")
			return nil, nil
		}
		log.Println("ErrNoRows2")
		return nil, err
	}
	return user, nil
}

// CreateUser creates user in Postgres
func (repo *UserPgRepo) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) 
	user.Password = string(hashedPassword)	

	_, err := repo.db.Model(user).Insert() 	//.TableExpr("users")
	if err != nil {
		return nil, errors.Wrap(err, "pg.CreateUser")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, _ := token.SignedString([]byte("secret"))
	user.Token = t

	user.Password = "" //delete password

	return user, nil
}

// CreateUser creates user in Postgres
func (repo *UserPgRepo) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	usercur := &model.User{} 
	err := repo.db.Model(usercur).Table("users").
			Where("users.email = ?", user.Email). 
			Select()
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usercur.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		log.Println("Invalid login credentials. Please try again")
		return nil, errors.Wrap(err, "Invalid login credentials. Please try again") 
	}

	usercur.Password = ""

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = usercur.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, _ := token.SignedString([]byte("secret"))
	usercur.Token = t

	return usercur, nil
}

