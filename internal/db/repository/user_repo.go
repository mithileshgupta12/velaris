package repository

import (
	"errors"

	"github.com/mithileshgupta12/velaris/internal/db/models"
	"xorm.io/xorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserCreationFailed = errors.New("failed to create user")
)

type UserRepository interface {
	CreateUser(args *CreateUserArgs) error
	GetUserById(userId int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	engine *xorm.Engine
}

func NewUserRepository(engine *xorm.Engine) UserRepository {
	return &userRepository{engine}
}

type CreateUserArgs struct {
	Name     string
	Email    string
	Password string
}

func (ur *userRepository) CreateUser(args *CreateUserArgs) error {
	user := &models.User{
		Name:     args.Name,
		Email:    args.Email,
		Password: args.Password,
	}

	affected, err := ur.engine.Insert(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserCreationFailed
	}

	return nil
}

func (ur *userRepository) GetUserById(userId int64) (*models.User, error) {
	user := new(models.User)

	has, err := ur.engine.
		Alias("u").
		Where("u.id = ?", userId).
		Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (ur *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)

	has, err := ur.engine.
		Alias("u").
		Where("u.email = ?", email).
		Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrUserNotFound
	}

	return user, nil
}
