package repository

import "xorm.io/xorm"

type Repository struct {
	UserRepository
	BoardRepository
	ListRepository
}

func NewRepository(engine *xorm.Engine) *Repository {
	userRepository := NewUserRepository(engine)
	boardRepository := NewBoardRepository(engine)
	listRepository := NewListRepository(engine)

	return &Repository{
		UserRepository:  userRepository,
		BoardRepository: boardRepository,
		ListRepository:  listRepository,
	}
}
