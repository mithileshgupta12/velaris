package db

import (
	"context"

	"github.com/mithileshgupta12/velaris/internal/db/repository"
)

type MockQueries struct {
	CreateBoardFunc  func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error)
	DeleteBoardFunc  func(ctx context.Context, id int32) (int64, error)
	GetAllBoardsFunc func(ctx context.Context) ([]repository.Board, error)
	GetBoardByIdFunc func(ctx context.Context, id int32) (repository.Board, error)
}

func (mq *MockQueries) CreateBoard(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
	if mq.CreateBoardFunc != nil {
		return mq.CreateBoardFunc(ctx, arg)
	}

	return repository.Board{}, nil
}

func (mq *MockQueries) DeleteBoard(ctx context.Context, id int32) (int64, error) {
	if mq.DeleteBoardFunc != nil {
		return mq.DeleteBoardFunc(ctx, id)
	}

	return 0, nil
}

func (mq *MockQueries) GetAllBoards(ctx context.Context) ([]repository.Board, error) {
	if mq.GetAllBoardsFunc != nil {
		return mq.GetAllBoardsFunc(ctx)
	}

	return []repository.Board{}, nil
}

func (mq *MockQueries) GetBoardById(ctx context.Context, id int32) (repository.Board, error) {
	if mq.GetBoardByIdFunc != nil {
		return mq.GetBoardByIdFunc(ctx, id)
	}

	return repository.Board{}, nil
}
