-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes
(
    mission_id INTEGER   NOT NULL,
    target_id  INTEGER   NOT NULL,

    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (mission_id, target_id) REFERENCES targets (mission_id, id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
-- +goose StatementEnd
