package service

import (
	"context"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

type CatBreedChecker interface {
	CheckBreed(ctx context.Context, breed string) (formattedBreed string, err error)
}

type CatsRepository interface {
	Create(ctx context.Context, params dto.CreateCatParams) (catID int, err error)
	Delete(ctx context.Context, catID int) (err error)
	Update(ctx context.Context, params dto.UpdateCatParams) (err error)
	One(ctx context.Context, catID int) (*models.Cat, error)
	All(ctx context.Context) ([]*models.Cat, error)
}

type MissionsRepository interface {
	Create(ctx context.Context) (missionID int, err error)
	Delete(ctx context.Context, missionID int) (err error)
	Update(ctx context.Context, params dto.UpdateMissionParams) (err error)
	One(ctx context.Context, missionID int) (*models.Mission, error)
	All(ctx context.Context, params dto.GetMissionsParams) ([]*models.Mission, error)
}

type TargetsRepository interface {
	Create(ctx context.Context, missionID int, lastTargetID int, targets []dto.CreateTargetParams) error
	Delete(ctx context.Context, missionID int, targetID int) (err error)
	Update(ctx context.Context, missionID int, targetID int, completed *bool) (err error)
	All(ctx context.Context, missionID int) ([]*models.Target, error)
	One(ctx context.Context, missionID int, targetID int) (*models.Target, error)
}

type NotesRepository interface {
	Create(ctx context.Context, missionID int, targetID int, contents []string) error
	All(ctx context.Context, missionID int, targetID int) ([]*models.Note, error)
}
