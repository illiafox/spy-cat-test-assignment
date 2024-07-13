package postgres

import (
	"context"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/postgres/schema"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
	"github.com/illiafox/spy-cat-test-assignment/app/pkg/poolwrapper"
	"github.com/jackc/pgx/v5"
	"time"
)

type MissionsRepository struct {
	db *poolwrapper.Pool
}

func NewMissionsRepository(db *poolwrapper.Pool) *MissionsRepository {
	return &MissionsRepository{db: db}
}

func (r *MissionsRepository) Create(ctx context.Context) (missionID int, err error) {
	var query = "INSERT INTO missions DEFAULT VALUES RETURNING id"

	err = r.db.QueryRow(ctx, query).Scan(&missionID)
	if err != nil {
		return -1, apperrors.Internal(err).Wrap("create mission: pgx: query row").
			WithMetadata("query", query)
	}

	return missionID, nil
}

func (r *MissionsRepository) Delete(ctx context.Context, missionID int) (err error) {
	query := "DELETE FROM missions WHERE id = $1"

	res, err := r.db.Exec(ctx, query, missionID)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("mission_id", missionID)
	}

	if res.RowsAffected() == 0 {
		return apperrors.MissionNotFound(missionID)
	}

	return nil
}

func (r *MissionsRepository) Update(ctx context.Context, params dto.UpdateMissionParams) (err error) {
	builder := sqlbuilder.Update("missions")
	builder.Set(builder.Assign("updated_at", time.Now()))

	if params.AssignedCatID != nil {
		builder.SetMore(builder.Assign("assigned_cat_id", *params.AssignedCatID))
	}

	if params.IsCompleted != nil {
		builder.SetMore(builder.Assign("is_completed", *params.IsCompleted))
	}

	query, args := builder.Where(builder.Equal("id", params.MissionID)).
		Build()

	//

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	if res.RowsAffected() == 0 {
		return apperrors.MissionNotFound(params.MissionID)
	}

	return nil
}

func (r *MissionsRepository) One(ctx context.Context, missionID int) (*models.Mission, error) {
	var mission schema.Mission

	const query = `SELECT id, assigned_cat_id, is_completed, created_at, updated_at
		FROM missions WHERE id = $1`

	rows, err := r.db.Query(ctx, query, missionID)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query mission").
			WithMetadata("query", query).
			WithMetadata("mission_id", missionID)
	}

	mission, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Mission])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.MissionNotFound(missionID)
		}

		return nil, apperrors.Internal(err).Wrap("pgx.CollectOneRow")
	}

	return mission.ToModel(), nil
}

func (r *MissionsRepository) All(ctx context.Context, params dto.GetMissionsParams) ([]*models.Mission, error) {
	var schemaMissions []schema.Mission

	builder := sqlbuilder.Select("id", "assigned_cat_id", "is_completed", "created_at", "updated_at").
		From("missions")

	if params.CatID != nil {
		builder.Where(builder.Equal("assigned_cat_id", *params.CatID))
	}

	query, args := builder.Build()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	schemaMissions, err = pgx.CollectRows(rows, pgx.RowToStructByName[schema.Mission])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.Internal(err).Wrap("pgx.CollectRows")
	}

	missions := make([]*models.Mission, len(schemaMissions))
	for i := range schemaMissions {
		missions[i] = schemaMissions[i].ToModel()
	}

	return missions, nil
}
