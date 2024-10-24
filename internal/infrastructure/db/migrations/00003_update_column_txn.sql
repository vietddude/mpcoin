-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE transactions ALTER COLUMN gas_price TYPE VARCHAR(255);
ALTER TABLE transactions ALTER COLUMN gas_limit TYPE VARCHAR(255);
ALTER TABLE transactions DROP COLUMN from_address;
ALTER TABLE chains ADD COLUMN explorer_url VARCHAR(255);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
ALTER TABLE transactions ALTER COLUMN gas_price TYPE NUMERIC;
ALTER TABLE transactions ALTER COLUMN gas_limit TYPE NUMERIC;
ALTER TABLE transactions ADD COLUMN from_address VARCHAR(255);
ALTER TABLE chains DROP COLUMN explorer_url;
