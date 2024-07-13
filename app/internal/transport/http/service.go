package http

import (
	"context"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
)

type Service interface {
	AddCat(ctx context.Context, params dto.CreateCatParams) (catID int, err error)
	GetCats(ctx context.Context) ([]*models.Cat, error)
	GetCatByID(ctx context.Context, catID int) (*models.Cat, error)
	UpdateCatByID(ctx context.Context, params dto.UpdateCatParams) error
	DeleteCatByID(ctx context.Context, catID int) error
	GetMissions(ctx context.Context, params dto.GetMissionsParams) ([]*models.Mission, error)
	CreateMission(ctx context.Context, params dto.CreateMissionParams) (missionID int, err error)
	AddMissionTargets(ctx context.Context, missionID int, newTargets []dto.CreateTargetParams) (err error)
	GetMissionByID(ctx context.Context, missionID int) (out *models.MissionFull, err error)
	UpdateMissionByID(ctx context.Context, params dto.UpdateMissionParams) (err error)
	DeleteMissionByID(ctx context.Context, missionID int) (err error)
	GetTargetsByMissionID(ctx context.Context, missionID int) (out []*models.TargetFull, err error)
	GetTargetByID(ctx context.Context, missionID int, targetID int) (out *models.TargetFull, err error)
	CompleteTargetByID(ctx context.Context, missionID int, targetID int) (err error)
	DeleteTargetByID(ctx context.Context, missionID int, targetID int) (err error)
	AddTargetNote(ctx context.Context, missionID int, targetID int, contents []string) (err error)
}
