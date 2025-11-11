package repository

import (
	"errors"

	"github.com/mithileshgupta12/velaris/internal/db/models"
	"xorm.io/xorm"
)

var (
	ErrBoardNotFound       = errors.New("board not found")
	ErrBoardCreationFailed = errors.New("failed to create board")
)

type BoardRepository interface {
	GetAllBoardsByUserId(userId int64) ([]*models.Board, error)
	CreateBoard(args *CreateBoardArgs) (*models.Board, error)
	GetBoardById(args *GetBoardByIdArgs) (*models.Board, error)
	UpdateBoardById(args *UpdateBoardByIdArgs) (*models.Board, error)
	DeleteBoardById(args *DeleteBoardByIdArgs) error
}

type boardRepository struct {
	engine *xorm.Engine
}

func NewBoardRepository(engine *xorm.Engine) *boardRepository {
	return &boardRepository{engine}
}

func (br *boardRepository) GetAllBoardsByUserId(userId int64) ([]*models.Board, error) {
	boards := []*models.Board{}

	err := br.engine.
		Alias("b").
		Where("b.userId = ?", userId).
		Find(&boards)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

type CreateBoardArgs struct {
	Name        string
	Description *string
	UserId      int64
}

func (br *boardRepository) CreateBoard(args *CreateBoardArgs) (*models.Board, error) {
	board := &models.Board{
		Name:        args.Name,
		Description: args.Description,
		UserId:      args.UserId,
	}

	affected, err := br.engine.
		Insert(board)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrBoardCreationFailed
	}

	return board, nil
}

type GetBoardByIdArgs struct {
	Id int64
}

func (br *boardRepository) GetBoardById(args *GetBoardByIdArgs) (*models.Board, error) {
	board := new(models.Board)

	has, err := br.engine.
		Alias("b").
		Where("b.id = ?", args.Id).
		Get(board)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrBoardNotFound
	}

	return board, nil
}

type UpdateBoardByIdArgs struct {
	Id          int64
	Name        string
	Description *string
}

func (br *boardRepository) UpdateBoardById(args *UpdateBoardByIdArgs) (*models.Board, error) {
	board := &models.Board{
		Name:        args.Name,
		Description: args.Description,
	}

	affected, err := br.engine.
		Alias("b").
		Where("b.id = ?", args.Id).
		Update(board)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, ErrBoardNotFound
	}

	return board, nil
}

type DeleteBoardByIdArgs struct {
	Id int64
}

func (br *boardRepository) DeleteBoardById(args *DeleteBoardByIdArgs) error {
	board := &models.Board{
		Id: args.Id,
	}

	affected, err := br.engine.
		Delete(board)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrBoardNotFound
	}

	return nil
}
