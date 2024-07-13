-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cats
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(100) NOT NULL,
    experience_years SMALLINT     NOT NULL,
    breed            VARCHAR(100) NOT NULL,
    salary           INTEGER      NOT NULL,

    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP    NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cats;
-- +goose StatementEnd
