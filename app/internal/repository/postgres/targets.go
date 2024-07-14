package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/postgres/schema"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
	"github.com/illiafox/spy-cat-test-assignment/app/pkg/poolwrapper"
	"github.com/jackc/pgx/v5"
)

type TargetsRepository struct {
	db *poolwrapper.Pool
}

func NewTargetsRepository(db *poolwrapper.Pool) *TargetsRepository {
	return &TargetsRepository{db: db}
}

func (r *TargetsRepository) Create(ctx context.Context, missionID int, lastTargetID int, targets []dto.CreateTargetParams) error {
	builder := sqlbuilder.InsertInto("targets").Cols("id", "mission_id", "name", "country")
	for _, target := range targets {
		lastTargetID++
		builder.Values(lastTargetID, missionID, target.Name, target.Country)
	}

	query, args := builder.Build()

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err).Wrap("create targets").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	return nil
}

func (r *TargetsRepository) Delete(ctx context.Context, missionID, targetID int) (err error) {
	query := "DELETE FROM targets WHERE mission_id = $1 AND id = $2"

	res, err := r.db.Exec(ctx, query, missionID, targetID)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("mission_id", missionID).
			WithMetadata("target_id", targetID)
	}

	if res.RowsAffected() == 0 {
		return apperrors.TargetNotFound(targetID)
	}

	return nil
}

func (r *TargetsRepository) Update(ctx context.Context, missionID, targetID int, completed *bool) (err error) {
	builder := sqlbuilder.Update("targets")
	builder.Set(builder.Assign("updated_at", time.Now()))

	if completed != nil {
		builder.SetMore(builder.Assign("is_completed", *completed))
	}

	query, args := builder.Where(
		builder.Equal("mission_id", missionID),
		builder.Equal("id", targetID),
	).Build()

	//

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	if res.RowsAffected() == 0 {
		return apperrors.TargetNotFound(targetID)
	}

	return nil
}

func (r *TargetsRepository) All(ctx context.Context, missionID int) ([]*models.Target, error) {
	var schemaTargets []schema.Target

	builder := sqlbuilder.Select("id", "mission_id", "is_completed", "name", "country", "created_at", "updated_at").
		From("targets").OrderBy("created_at").Desc()
	builder.Where(builder.Equal("mission_id", missionID))

	query, args := builder.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	schemaTargets, err = pgx.CollectRows(rows, pgx.RowToStructByName[schema.Target])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.Internal(err).Wrap("pgx.CollectRows")
	}

	targets := make([]*models.Target, len(schemaTargets))
	for i := range schemaTargets {
		targets[i] = schemaTargets[i].ToModel()
	}

	return targets, nil
}

func (r *TargetsRepository) One(ctx context.Context, missionID int, targetID int) (*models.Target, error) {
	var target schema.Target

	const query = `SELECT id, mission_id, is_completed, name, country, created_at, updated_at
		FROM targets WHERE mission_id = $1 AND id = $2 `

	rows, err := r.db.Query(ctx, query, missionID, targetID)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query target").
			WithMetadata("query", query).
			WithMetadata("mission_id", targetID)
	}

	target, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Target])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.TargetNotFound(targetID)
		}

		return nil, apperrors.Internal(err).Wrap("pgx.CollectOneRow")
	}

	return target.ToModel(), nil
}
