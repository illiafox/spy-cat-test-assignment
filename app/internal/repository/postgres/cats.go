package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/repository/postgres/schema"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
	"github.com/illiafox/spy-cat-test-assignment/app/pkg/poolwrapper"
	"github.com/jackc/pgx/v5"
)

type CatsRepository struct {
	db *poolwrapper.Pool
}

func NewCatsRepository(db *poolwrapper.Pool) *CatsRepository {
	return &CatsRepository{db: db}
}

func (r *CatsRepository) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	return r.db.TxFunc(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return f(ctx)
	})
}

func (r *CatsRepository) Create(ctx context.Context, params dto.CreateCatParams) (catID int, err error) {
	const query = `INSERT INTO cats(name, breed, experience_years, salary) VALUES ($1, $2, $3, $4)
	RETURNING id`
	args := []any{params.Name, params.Breed, params.Experience, params.Salary}

	err = r.db.QueryRow(ctx, query, args...).Scan(&catID)
	if err != nil {
		return -1, apperrors.Internal(err).Wrap("pgx: query row").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	return catID, nil
}

func (r *CatsRepository) Delete(ctx context.Context, catID int) (err error) {
	const query = "DELETE FROM cats WHERE id = $1"

	res, err := r.db.Exec(ctx, query, catID)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("cat_id", catID)
	}

	if res.RowsAffected() == 0 {
		return apperrors.CatNotFound(catID)
	}

	return nil
}

func (r *CatsRepository) Update(ctx context.Context, params dto.UpdateCatParams) (err error) {
	builder := sqlbuilder.Update("cats")
	builder.Set(builder.Assign("updated_at", time.Now()))

	if params.Name != nil {
		builder.SetMore(builder.Assign("name", *params.Name))
	}

	if params.ExperienceYears != nil {
		builder.SetMore(builder.Assign("experience_years", *params.ExperienceYears))
	}

	if params.Salary != nil {
		builder.SetMore(builder.Assign("salary", *params.Salary))
	}

	query, args := builder.Where(builder.Equal("id", params.CatID)).
		Build()

	//

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err).Wrap("pgx: exec").
			WithMetadata("query", query).
			WithMetadata("args", args)
	}

	if res.RowsAffected() == 0 {
		return apperrors.CatNotFound(params.CatID)
	}

	return nil
}

func (r *CatsRepository) One(ctx context.Context, catID int) (*models.Cat, error) {
	var cat schema.Cat

	const query = `SELECT id, name, experience_years, breed, salary, created_at, updated_at
		FROM cats WHERE id = $1`

	rows, err := r.db.Query(ctx, query, catID)
	if err != nil {
		return nil, apperrors.Internal(err).Wrap("pgx: query").
			WithMetadata("query", query).
			WithMetadata("cat_id", catID)
	}

	cat, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Cat])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.CatNotFound(catID)
		}

		return nil, apperrors.Internal(err).Wrap("pgx.CollectOneRow")
	}

	return cat.ToModel(), nil
}

func (r *CatsRepository) All(ctx context.Context) ([]*models.Cat, error) {
	var schemaCats []schema.Cat

	const query = `SELECT id, name, experience_years, breed, salary, created_at, updated_at FROM cats`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, apperrors.Internal(err).Wrap("pgx: query").
			WithMetadata("query", query)
	}

	schemaCats, err = pgx.CollectRows(rows, pgx.RowToStructByName[schema.Cat])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.Internal(err).Wrap("pgx.CollectRows")
	}

	cats := make([]*models.Cat, len(schemaCats))
	for i := range schemaCats {
		cats[i] = schemaCats[i].ToModel()
	}

	return cats, nil
}
