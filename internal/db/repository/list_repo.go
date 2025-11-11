package repository

import (
	"errors"

	"github.com/mithileshgupta12/velaris/internal/db/models"
	"xorm.io/xorm"
)

var (
	ErrListNotFound = errors.New("list not found")
)

type ListRepository interface {
	GetAllListsByBoardId(args *GetAllListsByBoardIdArgs) ([]*models.List, error)
	DeleteListById(args *DeleteListByIdArgs) error
}

type listRepository struct {
	engine *xorm.Engine
}

func NewListRepository(engine *xorm.Engine) ListRepository {
	return &listRepository{engine}
}

type GetAllListsByBoardIdArgs struct {
	BoardId int64
}

func (lr *listRepository) GetAllListsByBoardId(args *GetAllListsByBoardIdArgs) ([]*models.List, error) {
	lists := []*models.List{}

	err := lr.engine.
		Alias("l").
		Where("l.board_id = ?", args.BoardId).
		Find(&lists)
	if err != nil {
		return nil, err
	}

	return lists, nil
}

type DeleteListByIdArgs struct {
	ListId int64
}

func (lr *listRepository) DeleteListById(args *DeleteListByIdArgs) error {
	list := &models.List{
		Id: args.ListId,
	}

	affected, err := lr.engine.
		Delete(list)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrListNotFound
	}

	return nil
}
