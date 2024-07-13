-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS targets
(
    mission_id   INTEGER      NOT NULL,
    id           INTEGER      NOT NULL,

    is_completed BOOLEAN               DEFAULT FALSE,

    name         VARCHAR(100) NOT NULL,
    country      CHAR(2)      NOT NULL, -- ISO 3166-1 alpha-2

    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT NOW(),

    PRIMARY KEY (mission_id, id),
    FOREIGN KEY (mission_id) REFERENCES missions (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS targets;
-- +goose StatementEnd
