package repository

import (
	"github.com/mithileshgupta12/velaris/internal/db/models"
	"xorm.io/xorm"
)

type ListRepository interface {
	GetAllListsByBoardId(args *GetAllListsByBoardIdArgs) ([]*models.List, error)
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
