package postgres

import (
	"context"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/postgres/schema"
	"github.com/illiafox/spy-cat-test-assignment/app/pkg/poolwrapper"
	"github.com/jackc/pgx/v5"
)

type NotesRepository struct {
	db *poolwrapper.Pool
}

func NewNotesRepository(db *poolwrapper.Pool) *NotesRepository {
	return &NotesRepository{db: db}
}

func (r *NotesRepository) Create(ctx context.Context, missionID, targetID int, contents []string) error {
	builder := sqlbuilder.InsertInto("notes").Cols("mission_id", "target_id", "content")
	for _, content := range contents {
		builder.Values(missionID, targetID, content)
	}

	query, args := builder.Build()

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err).Wrap("create notes: pgx: exec").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	return nil
}

func (r *NotesRepository) All(ctx context.Context, missionID, targetID int) ([]*models.Note, error) {
	var schemaNotes []schema.Note

	builder := sqlbuilder.Select("mission_id", "target_id", "content", "created_at").
		From("notes").OrderBy("created_at").Desc()

	query, args := builder.Where(
		builder.Equal("mission_id", missionID),
		builder.Equal("target_id", targetID),
	).Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	schemaNotes, err = pgx.CollectRows(rows, pgx.RowToStructByName[schema.Note])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.Internal(err).Wrap("pgx.CollectRows")
	}

	notes := make([]*models.Note, len(schemaNotes))
	for i := range schemaNotes {
		notes[i] = schemaNotes[i].ToModel()
	}

	return notes, nil
}
