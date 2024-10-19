-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE transactions
ALTER COLUMN amount TYPE VARCHAR(255);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE transactions
ALTER COLUMN amount TYPE NUMERIC(78, 0);
