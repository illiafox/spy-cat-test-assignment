-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS missions
(
    id              SERIAL PRIMARY KEY,
    assigned_cat_id INTEGER,
    is_completed    BOOLEAN      NOT NULL DEFAULT false,

    created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP    NOT NULL DEFAULT NOW(),

    FOREIGN KEY (assigned_cat_id) REFERENCES cats (id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS missions;
-- +goose StatementEnd
