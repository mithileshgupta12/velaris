package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
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

func TestBoardHandler_Index(t *testing.T) {
	now := time.Now()

	boards := []repository.Board{
		{
			ID:   1,
			Name: "lorem",
			Description: pgtype.Text{
				String: "lorem ipsum dolor",
				Valid:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		Name        string
		Endpoint    string
		Method      string
		MockQueries *MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 200 and list of boards when repository returns boards successfully",
			Endpoint: "/boards",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetAllBoardsFunc: func(ctx context.Context) ([]repository.Board, error) {
					return boards, nil
				},
			},
			StatusCode: http.StatusOK,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boards,
			},
			Success: true,
		},
		{
			Name:     "must return 200 and empty array when repository returns no boards",
			Endpoint: "/boards",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetAllBoardsFunc: func(ctx context.Context) ([]repository.Board, error) {
					return []repository.Board{}, nil
				},
			},
			StatusCode: http.StatusOK,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    []repository.Board{},
			},
			Success: true,
		},
		{
			Name:     "must return 500 and error message when repository returns error",
			Endpoint: "/boards",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetAllBoardsFunc: func(ctx context.Context) ([]repository.Board, error) {
					return []repository.Board{}, errors.New("some error")
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "internal server error",
				},
			},
			Success: false,
		},
	}

	lgr := logger.NewLogger(logger.FormatJSON)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			boardHandler := NewBoardHandler(test.MockQueries, lgr)

			r := httptest.NewRequest(test.Method, test.Endpoint, nil)
			rr := httptest.NewRecorder()

			boardHandler.Index(rr, r)

			if test.StatusCode != rr.Code {
				t.Errorf("want status code %d got %d", test.StatusCode, rr.Code)
			}

			wantHeader, gotHeader := "application/json", rr.Header().Get("Content-Type")

			if gotHeader != wantHeader {
				t.Errorf("want %s got %s", wantHeader, gotHeader)
			}

			if !test.Success {
				var gotResponse helper.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &gotResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if !reflect.DeepEqual(test.Response, gotResponse) {
					t.Errorf("wanted %+v got %+v", test.Response, gotResponse)
				}
			} else {
				expectedJSON, err := json.Marshal(test.Response)
				if err != nil {
					t.Errorf("failed to marshal json %v", err)
				}

				var expectedResponse, gotResponse helper.SuccessResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &gotResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if err := json.Unmarshal(expectedJSON, &expectedResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if !reflect.DeepEqual(expectedResponse, gotResponse) {
					t.Errorf("wanted %+v got %+v", expectedResponse, gotResponse)
				}
			}
		})
	}
}

func TestBoardHandler_Show(t *testing.T) {
	now := time.Now()

	board := repository.Board{
		ID:   1,
		Name: "lorem",
		Description: pgtype.Text{
			String: "lorem ipsum dolor",
			Valid:  true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		Name        string
		Endpoint    string
		ID          string
		Method      string
		MockQueries *MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 200 and board when repository returns board successfully",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetBoardByIdFunc: func(ctx context.Context, id int32) (repository.Board, error) {
					return board, nil
				},
			},
			StatusCode: http.StatusOK,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    board,
			},
			Success: true,
		},
		{
			Name:     "must return 404 and error message when board not found",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetBoardByIdFunc: func(ctx context.Context, id int32) (repository.Board, error) {
					return repository.Board{}, pgx.ErrNoRows
				},
			},
			StatusCode: http.StatusNotFound,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "board not found",
				},
			},
			Success: false,
		},
		{
			Name:        "must return 400 and error message when id parameter is invalid",
			Endpoint:    "/boards",
			ID:          "abc",
			Method:      http.MethodGet,
			MockQueries: &MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "invalid board id",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 500 and error message when repository returns unexpected error",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodGet,
			MockQueries: &MockQueries{
				GetBoardByIdFunc: func(ctx context.Context, id int32) (repository.Board, error) {
					return repository.Board{}, errors.New("some error")
				},
			},
			StatusCode: http.StatusInternalServerError,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "internal server error",
				},
			},
			Success: false,
		},
	}

	lgr := logger.NewLogger(logger.FormatJSON)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			boardHandler := NewBoardHandler(test.MockQueries, lgr)

			router := chi.NewRouter()

			router.Get("/boards/{id}", boardHandler.Show)

			r := httptest.NewRequest(test.Method, fmt.Sprintf("%s/%s", test.Endpoint, test.ID), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, r)

			if test.StatusCode != rr.Code {
				t.Errorf("want status code %d got %d", test.StatusCode, rr.Code)
			}

			wantHeader, gotHeader := "application/json", rr.Header().Get("Content-Type")

			if gotHeader != wantHeader {
				t.Errorf("want %s got %s", wantHeader, gotHeader)
			}

			if !test.Success {
				var gotResponse helper.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &gotResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if !reflect.DeepEqual(test.Response, gotResponse) {
					t.Errorf("wanted %+v got %+v", test.Response, gotResponse)
				}
			} else {
				expectedJSON, err := json.Marshal(test.Response)
				if err != nil {
					t.Errorf("failed to marshal json %v", err)
				}

				var expectedResponse, gotResponse helper.SuccessResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &gotResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if err := json.Unmarshal(expectedJSON, &expectedResponse); err != nil {
					t.Errorf("failed to unmarshal json %v", err)
				}

				if !reflect.DeepEqual(expectedResponse, gotResponse) {
					t.Errorf("wanted %+v got %+v", expectedResponse, gotResponse)
				}
			}
		})
	}
}
