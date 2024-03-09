-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS PUBLIC.EMPLOYEE (
    EMPL_ID SERIAL PRIMARY KEY,
    NAME VARCHAR
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT
    'DROP TABLE Employee';

-- +goose StatementEnd