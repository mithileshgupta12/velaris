-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS IDX_boards_user_id ON boards (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS IDX_boards_user_id;
-- +goose StatementEnd
