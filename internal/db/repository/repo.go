package repository

import "xorm.io/xorm"

type Repository struct {
	UserRepository
	BoardRepository
}

func NewRepository(engine *xorm.Engine) *Repository {
	userRepository := NewUserRepository(engine)
	boardRepository := NewBoardRepository(engine)

	return &Repository{UserRepository: userRepository, BoardRepository: boardRepository}
}
