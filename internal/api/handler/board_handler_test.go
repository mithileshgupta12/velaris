package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/db/repository"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func TestBoardHandler_Index(t *testing.T) {
	now := time.Now()

	boards := []repository.Board{
		{
			ID:   1,
			Name: "Sprint Board",
			Description: pgtype.Text{
				String: "Q4 2024 Sprint Planning",
				Valid:  true,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:   2,
			Name: "Backlog",
			Description: pgtype.Text{
				String: "Product backlog items",
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
		MockQueries *db.MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 200 and list of boards when repository returns boards successfully",
			Endpoint: "/boards",
			Method:   http.MethodGet,
			MockQueries: &db.MockQueries{
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
			MockQueries: &db.MockQueries{
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
			MockQueries: &db.MockQueries{
				GetAllBoardsFunc: func(ctx context.Context) ([]repository.Board, error) {
					return []repository.Board{}, errors.New("database connection failed")
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

	lgr := logger.NewTestLogger(logger.FormatJSON)

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
		Name: "Sprint Board",
		Description: pgtype.Text{
			String: "Q4 2024 Sprint Planning",
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
		MockQueries *db.MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 200 and board when repository returns board successfully",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodGet,
			MockQueries: &db.MockQueries{
				GetBoardByIdFunc: func(ctx context.Context, id int32) (repository.Board, error) {
					if id != 1 {
						t.Errorf("expected id 1, got %d", id)
					}
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
			ID:       "999",
			Method:   http.MethodGet,
			MockQueries: &db.MockQueries{
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
			MockQueries: &db.MockQueries{},
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
			Name:        "must return 400 and error message when id parameter is negative",
			Endpoint:    "/boards",
			ID:          "-1",
			Method:      http.MethodGet,
			MockQueries: &db.MockQueries{},
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
			MockQueries: &db.MockQueries{
				GetBoardByIdFunc: func(ctx context.Context, id int32) (repository.Board, error) {
					return repository.Board{}, errors.New("database connection failed")
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

	lgr := logger.NewTestLogger(logger.FormatJSON)

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

func TestBoardHandler_Store(t *testing.T) {
	now := time.Now()

	board := repository.Board{
		ID:   1,
		Name: "Sprint Board",
		Description: pgtype.Text{
			String: "Q4 2024 Sprint Planning",
			Valid:  true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	boardWithoutDescription := repository.Board{
		ID:   1,
		Name: "Sprint Board",
		Description: pgtype.Text{
			String: "",
			Valid:  false,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		Name        string
		Endpoint    string
		Method      string
		RequestBody any
		MockQueries *db.MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 201 and board when repository creates board successfully with description",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]string{
				"name":        "Sprint Board",
				"description": "Q4 2024 Sprint Planning",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Name != "Sprint Board" {
						t.Errorf("expected name 'Sprint Board', got '%s'", arg.Name)
					}
					if !arg.Description.Valid || arg.Description.String != "Q4 2024 Sprint Planning" {
						t.Errorf("expected valid description 'Q4 2024 Sprint Planning', got valid=%v, value='%s'",
							arg.Description.Valid, arg.Description.String)
					}
					return board, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    board,
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when repository creates board successfully without description",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name": "Sprint Board",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Name != "Sprint Board" {
						t.Errorf("expected name 'Sprint Board', got '%s'", arg.Name)
					}
					if arg.Description.Valid {
						t.Errorf("expected invalid description (no description provided), got valid=%v", arg.Description.Valid)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when description is empty string",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name":        "Sprint Board",
				"description": "",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Description.Valid {
						t.Errorf("expected invalid description for empty string, got valid=%v", arg.Description.Valid)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
		{
			Name:     "must return 400 and error message when name is empty string",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name": "",
			},
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "name is a required field",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 400 and error message when name contains only whitespace",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name": "     ",
			},
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "name is a required field",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 201 and board when name has leading/trailing whitespace (should trim)",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name": "  Sprint Board  ",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Name != "Sprint Board" {
						t.Errorf("expected trimmed name 'Sprint Board', got '%s'", arg.Name)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when description has leading/trailing whitespace (should trim)",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name":        "Sprint Board",
				"description": "  Q4 2024 Sprint Planning  ",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Description.String != "Q4 2024 Sprint Planning" {
						t.Errorf("expected trimmed description 'Q4 2024 Sprint Planning', got '%s'", arg.Description.String)
					}
					return board, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    board,
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when description contains only whitespace (should be empty)",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]any{
				"name":        "Sprint Board",
				"description": "     ",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Description.Valid {
						t.Errorf("expected invalid description for whitespace-only input, got valid=%v", arg.Description.Valid)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
		{
			Name:        "must return 400 and error message when request body is invalid JSON",
			Endpoint:    "/boards",
			Method:      http.MethodPost,
			RequestBody: "invalid json{",
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "invalid request",
				},
			},
			Success: false,
		},
		{
			Name:        "must return 400 and error message when request body is empty",
			Endpoint:    "/boards",
			Method:      http.MethodPost,
			RequestBody: "",
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "invalid request",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 500 and error message when repository returns error",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]string{
				"name":        "Sprint Board",
				"description": "Q4 2024 Sprint Planning",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					return repository.Board{}, errors.New("database connection failed")
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
		{
			Name:     "must return 201 and board when name contains special characters",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]string{
				"name":        "Sprint-Board_2024!",
				"description": "Special chars test",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Name != "Sprint-Board_2024!" {
						t.Errorf("expected name with special chars, got '%s'", arg.Name)
					}
					return repository.Board{
						ID:   1,
						Name: "Sprint-Board_2024!",
						Description: pgtype.Text{
							String: "Special chars test",
							Valid:  true,
						},
					}, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data: repository.Board{
					ID:   1,
					Name: "Sprint-Board_2024!",
					Description: pgtype.Text{
						String: "Special chars test",
						Valid:  true,
					},
				},
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when name contains unicode characters",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]string{
				"name":        "Sprint Board 🚀",
				"description": "Unicode test",
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if arg.Name != "Sprint Board 🚀" {
						t.Errorf("expected name with unicode, got '%s'", arg.Name)
					}
					return repository.Board{
						ID:   1,
						Name: "Sprint Board 🚀",
						Description: pgtype.Text{
							String: "Unicode test",
							Valid:  true,
						},
					}, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data: repository.Board{
					ID:   1,
					Name: "Sprint Board 🚀",
					Description: pgtype.Text{
						String: "Unicode test",
						Valid:  true,
					},
				},
			},
			Success: true,
		},
		{
			Name:     "must return 201 and board when description is very long",
			Endpoint: "/boards",
			Method:   http.MethodPost,
			RequestBody: map[string]string{
				"name":        "Sprint Board",
				"description": strings.Repeat("A", 1000),
			},
			MockQueries: &db.MockQueries{
				CreateBoardFunc: func(ctx context.Context, arg repository.CreateBoardParams) (repository.Board, error) {
					if len(arg.Description.String) != 1000 {
						t.Errorf("expected description length 1000, got %d", len(arg.Description.String))
					}
					return repository.Board{
						ID:   1,
						Name: "Sprint Board",
						Description: pgtype.Text{
							String: strings.Repeat("A", 1000),
							Valid:  true,
						},
					}, nil
				},
			},
			StatusCode: http.StatusCreated,
			Response: helper.SuccessResponse{
				Success: true,
				Data: repository.Board{
					ID:   1,
					Name: "Sprint Board",
					Description: pgtype.Text{
						String: strings.Repeat("A", 1000),
						Valid:  true,
					},
				},
			},
			Success: true,
		},
	}

	lgr := logger.NewTestLogger(logger.FormatJSON)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var requestBody io.Reader
			if test.RequestBody != nil {
				if strBody, ok := test.RequestBody.(string); ok {
					requestBody = strings.NewReader(strBody)
				} else {
					requestBytes, err := json.Marshal(test.RequestBody)
					if err != nil {
						t.Fatalf("failed to marshal json %v", err)
					}
					requestBody = bytes.NewReader(requestBytes)
				}
			}

			boardHandler := NewBoardHandler(test.MockQueries, lgr)

			r := httptest.NewRequest(test.Method, test.Endpoint, requestBody)
			r.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			boardHandler.Store(rr, r)

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

func TestBoardHandler_Update(t *testing.T) {
	now := time.Now()

	board := repository.Board{
		ID:   1,
		Name: "Updated Sprint Board",
		Description: pgtype.Text{
			String: "Updated description",
			Valid:  true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	boardWithoutDescription := repository.Board{
		ID:   1,
		Name: "Updated Sprint Board",
		Description: pgtype.Text{
			String: "",
			Valid:  false,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		Name        string
		Endpoint    string
		ID          string
		Method      string
		RequestBody any
		MockQueries *db.MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 200 and updated board when repository updates board successfully with description",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]string{
				"name":        "Updated Sprint Board",
				"description": "Updated description",
			},
			MockQueries: &db.MockQueries{
				UpdateBoardByIdFunc: func(ctx context.Context, arg repository.UpdateBoardByIdParams) (repository.Board, error) {
					if arg.ID != 1 {
						t.Errorf("expected id 1, got %d", arg.ID)
					}
					if arg.Name != "Updated Sprint Board" {
						t.Errorf("expected name 'Updated Sprint Board', got '%s'", arg.Name)
					}
					if !arg.Description.Valid || arg.Description.String != "Updated description" {
						t.Errorf("expected valid description 'Updated description', got valid=%v, value='%s'",
							arg.Description.Valid, arg.Description.String)
					}
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
			Name:     "must return 200 and updated board when repository updates board successfully without description",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]any{
				"name": "Updated Sprint Board",
			},
			MockQueries: &db.MockQueries{
				UpdateBoardByIdFunc: func(ctx context.Context, arg repository.UpdateBoardByIdParams) (repository.Board, error) {
					if arg.Description.Valid {
						t.Errorf("expected invalid description, got valid=%v", arg.Description.Valid)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusOK,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
		{
			Name:     "must return 404 and error message when board not found",
			Endpoint: "/boards",
			ID:       "999",
			Method:   http.MethodPut,
			RequestBody: map[string]string{
				"name":        "Updated Sprint Board",
				"description": "Updated description",
			},
			MockQueries: &db.MockQueries{
				UpdateBoardByIdFunc: func(ctx context.Context, arg repository.UpdateBoardByIdParams) (repository.Board, error) {
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
			Method:      http.MethodPut,
			RequestBody: map[string]string{"name": "Test"},
			MockQueries: &db.MockQueries{},
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
			Name:     "must return 400 and error message when name is empty string",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]any{
				"name": "",
			},
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "name is a required field",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 400 and error message when name contains only whitespace",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]any{
				"name": "     ",
			},
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "name is a required field",
				},
			},
			Success: false,
		},
		{
			Name:        "must return 400 and error message when request body is invalid JSON",
			Endpoint:    "/boards",
			ID:          "1",
			Method:      http.MethodPut,
			RequestBody: "invalid json{",
			MockQueries: &db.MockQueries{},
			StatusCode:  http.StatusBadRequest,
			Response: helper.ErrorResponse{
				Success: false,
				Error: helper.Error{
					Message: "invalid request",
				},
			},
			Success: false,
		},
		{
			Name:     "must return 500 and error message when repository returns unexpected error",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]string{
				"name":        "Updated Sprint Board",
				"description": "Updated description",
			},
			MockQueries: &db.MockQueries{
				UpdateBoardByIdFunc: func(ctx context.Context, arg repository.UpdateBoardByIdParams) (repository.Board, error) {
					return repository.Board{}, errors.New("database connection failed")
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
		{
			Name:     "must return 200 and board when name has leading/trailing whitespace (should trim)",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodPut,
			RequestBody: map[string]any{
				"name": "  Updated Sprint Board  ",
			},
			MockQueries: &db.MockQueries{
				UpdateBoardByIdFunc: func(ctx context.Context, arg repository.UpdateBoardByIdParams) (repository.Board, error) {
					if arg.Name != "Updated Sprint Board" {
						t.Errorf("expected trimmed name 'Updated Sprint Board', got '%s'", arg.Name)
					}
					return boardWithoutDescription, nil
				},
			},
			StatusCode: http.StatusOK,
			Response: helper.SuccessResponse{
				Success: true,
				Data:    boardWithoutDescription,
			},
			Success: true,
		},
	}

	lgr := logger.NewTestLogger(logger.FormatJSON)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var requestBody io.Reader
			if test.RequestBody != nil {
				if strBody, ok := test.RequestBody.(string); ok {
					requestBody = strings.NewReader(strBody)
				} else {
					requestBytes, err := json.Marshal(test.RequestBody)
					if err != nil {
						t.Fatalf("failed to marshal json %v", err)
					}
					requestBody = bytes.NewReader(requestBytes)
				}
			}

			boardHandler := NewBoardHandler(test.MockQueries, lgr)

			router := chi.NewRouter()
			router.Put("/boards/{id}", boardHandler.Update)

			r := httptest.NewRequest(test.Method, fmt.Sprintf("%s/%s", test.Endpoint, test.ID), requestBody)
			r.Header.Set("Content-Type", "application/json")
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

func TestBoardHandler_Destroy(t *testing.T) {
	tests := []struct {
		Name        string
		Endpoint    string
		ID          string
		Method      string
		MockQueries *db.MockQueries
		StatusCode  int
		Response    any
		Success     bool
	}{
		{
			Name:     "must return 204 when board is deleted successfully",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodDelete,
			MockQueries: &db.MockQueries{
				DeleteBoardByIdFunc: func(ctx context.Context, id int32) (int64, error) {
					if id != 1 {
						t.Errorf("expected id 1, got %d", id)
					}
					return 1, nil
				},
			},
			StatusCode: http.StatusNoContent,
			Response:   nil,
			Success:    true,
		},
		{
			Name:     "must return 404 and error message when board not found",
			Endpoint: "/boards",
			ID:       "999",
			Method:   http.MethodDelete,
			MockQueries: &db.MockQueries{
				DeleteBoardByIdFunc: func(ctx context.Context, id int32) (int64, error) {
					return 0, nil
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
			Method:      http.MethodDelete,
			MockQueries: &db.MockQueries{},
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
			Name:        "must return 400 and error message when id parameter is negative",
			Endpoint:    "/boards",
			ID:          "-1",
			Method:      http.MethodDelete,
			MockQueries: &db.MockQueries{},
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
			Name:     "must return 500 and error message when repository returns error",
			Endpoint: "/boards",
			ID:       "1",
			Method:   http.MethodDelete,
			MockQueries: &db.MockQueries{
				DeleteBoardByIdFunc: func(ctx context.Context, id int32) (int64, error) {
					return 0, errors.New("database connection failed")
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

	lgr := logger.NewTestLogger(logger.FormatJSON)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			boardHandler := NewBoardHandler(test.MockQueries, lgr)

			router := chi.NewRouter()
			router.Delete("/boards/{id}", boardHandler.Destroy)

			r := httptest.NewRequest(test.Method, fmt.Sprintf("%s/%s", test.Endpoint, test.ID), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, r)

			if test.StatusCode != rr.Code {
				t.Errorf("want status code %d got %d", test.StatusCode, rr.Code)
			}

			// For 204 No Content, there should be no body
			if test.StatusCode == http.StatusNoContent {
				if rr.Body.Len() > 0 {
					t.Errorf("expected empty body for 204, got %s", rr.Body.String())
				}
				return
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
			}
		})
	}
}
