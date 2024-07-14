package service

import (
	"context"
	"fmt"

	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
)

type Service struct {
	catBreedChecker    CatBreedChecker
	catsRepository     CatsRepository
	missionsRepository MissionsRepository
	targetsRepository  TargetsRepository
	notesRepository    NotesRepository

	transactor Transactor
}

func NewService(catBreedChecker CatBreedChecker, catsRepository CatsRepository, missionsRepository MissionsRepository, targetsRepository TargetsRepository, notesRepository NotesRepository, transactor Transactor) *Service {
	return &Service{catBreedChecker: catBreedChecker, catsRepository: catsRepository, missionsRepository: missionsRepository, targetsRepository: targetsRepository, notesRepository: notesRepository, transactor: transactor}
}

func (s Service) AddCat(ctx context.Context, params dto.CreateCatParams) (catID int, err error) {
	formattedBreed, err := s.catBreedChecker.CheckBreed(ctx, params.Breed)
	if err != nil {
		return -1, fmt.Errorf("check formattedBreed: %w", err)
	}

	params.Breed = formattedBreed
	return s.catsRepository.Create(ctx, params)
}

func (s Service) GetCats(ctx context.Context) ([]*models.Cat, error) {
	cats, err := s.catsRepository.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("cats repository: all: %w", err)
	}

	return cats, nil
}

func (s Service) GetCatByID(ctx context.Context, catID int) (*models.Cat, error) {
	cat, err := s.catsRepository.One(ctx, catID)
	if err != nil {
		return nil, fmt.Errorf("cats repository: one: %w", err)
	}

	return cat, nil
}

func (s Service) UpdateCatByID(ctx context.Context, params dto.UpdateCatParams) error {
	err := s.catsRepository.Update(ctx, params)
	if err != nil {
		return fmt.Errorf("cats repository: update: %w", err)
	}

	return nil
}

func (s Service) DeleteCatByID(ctx context.Context, catID int) error {
	err := s.catsRepository.Delete(ctx, catID)
	if err != nil {
		return fmt.Errorf("cats repository: delete: %w", err)
	}

	return nil
}

// Missions

func (s Service) GetMissions(ctx context.Context, params dto.GetMissionsParams) ([]*models.Mission, error) {
	missions, err := s.missionsRepository.All(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("missions repository: all: %w", err)
	}

	return missions, nil
}

func (s Service) CreateMission(ctx context.Context, params dto.CreateMissionParams) (missionID int, err error) {
	if targetsCount := len(params.Targets); targetsCount == 0 || targetsCount > 3 { // minimum: 1, maximum: 3
		return -1, apperrors.InvalidTargetsCount(targetsCount, 1, 3)
	}

	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		missionID, err = s.missionsRepository.Create(ctx)
		if err != nil {
			return fmt.Errorf("create mission: %w", err)
		}

		const lastTargetID = 0
		err = s.targetsRepository.Create(ctx, missionID, lastTargetID, params.Targets)
		if err != nil {
			return fmt.Errorf("create targets for mission %d: %w", missionID, err)
		}

		return nil
	})
	if err != nil {
		return -1, fmt.Errorf("within transaction: %w", err)
	}

	return missionID, nil
}

func (s Service) AddMissionTargets(ctx context.Context, missionID int, newTargets []dto.CreateTargetParams) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		mission, err := s.GetMissionByID(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		if mission.IsCompleted {
			return apperrors.MissionAlreadyCompleted(missionID)
		}

		targetsCount := len(mission.Targets) + len(newTargets)
		if targetsCount == 0 || targetsCount > 3 { // minimum: 1, maximum: 3
			return apperrors.InvalidTargetsCount(targetsCount, 1, 3).
				Wrap("too many targets")
		}

		err = s.targetsRepository.Create(ctx, missionID, len(mission.Targets), newTargets)
		if err != nil {
			return fmt.Errorf("create targets for mission %d: %w", missionID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}

func (s Service) GetMissionByID(ctx context.Context, missionID int) (out *models.MissionFull, err error) {
	out = new(models.MissionFull)

	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		out.Mission, err = s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission by id: %w", err)
		}

		out.Targets, err = s.targetsRepository.All(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get targets by mission id %d: %w", missionID, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("within transaction: %w", err)
	}

	return out, nil
}

func (s Service) UpdateMissionByID(ctx context.Context, params dto.UpdateMissionParams) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		mission, err := s.missionsRepository.One(ctx, params.MissionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", params.MissionID, err)
		}

		if mission.IsCompleted {
			return apperrors.MissionAlreadyCompleted(params.MissionID)
		}

		if params.IsCompleted != nil { // mission can be marked as completed only if all targets are completed
			targets, err := s.targetsRepository.All(ctx, params.MissionID)
			if err != nil {
				return fmt.Errorf("get targets: %w", err)
			}

			for _, t := range targets {
				if !t.IsCompleted {
					return apperrors.AllTargetsAreNotCompleted().Wrap("can't complete mission")
				}
			}
		}

		if params.AssignedCatID != nil {
			_, err := s.catsRepository.One(ctx, *params.AssignedCatID)
			if err != nil {
				return fmt.Errorf("get cat: %w", err)
			}
		}

		err = s.missionsRepository.Update(ctx, params)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("update mission %d: %w", params.MissionID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}

func (s Service) DeleteMissionByID(ctx context.Context, missionID int) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		mission, err := s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		if mission.AssignedCatID != nil {
			return apperrors.CatAlreadyAssigned(missionID).Wrap("can't delete mission")
		}

		if mission.IsCompleted {
			return apperrors.MissionAlreadyCompleted(missionID)
		}

		err = s.missionsRepository.Delete(ctx, missionID)
		if err != nil {
			return fmt.Errorf("delete mission %d: %w", missionID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}

// Targets

func (s Service) GetTargetsByMissionID(ctx context.Context, missionID int) (out []*models.TargetFull, err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		targets, err := s.targetsRepository.All(ctx, missionID)
		if err != nil {
			return fmt.Errorf("targets repo: all: %w", err)
		}

		out = make([]*models.TargetFull, len(targets))

		for i := range targets {
			out[i] = &models.TargetFull{
				Target: targets[i],
			}

			out[i].Notes, err = s.notesRepository.All(ctx, missionID, targets[i].ID)
			if err != nil {
				return fmt.Errorf("notes repo: get notes for target %d: %w", targets[i].ID, err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("within transaction: %w", err)
	}

	return out, nil
}

func (s Service) GetTargetByID(ctx context.Context, missionID, targetID int) (out *models.TargetFull, err error) {
	out = new(models.TargetFull)

	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err = s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		out.Target, err = s.targetsRepository.One(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("get target by id: %w", err)
		}

		if out.Target.MissionID != missionID {
			return apperrors.TargetNotFound(targetID).Wrap("mission")
		}

		out.Notes, err = s.notesRepository.All(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("get notes by target id %d: %w", targetID, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("within transaction: %w", err)
	}

	return out, nil
}

func (s Service) CompleteTargetByID(ctx context.Context, missionID, targetID int) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err = s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		target, err := s.targetsRepository.One(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("get target by id: %w", err)
		}

		if target.IsCompleted {
			return apperrors.TargetAlreadyCompleted(targetID)
		}

		isCompleted := true

		err = s.targetsRepository.Update(ctx, missionID, targetID, &isCompleted)
		if err != nil {
			return fmt.Errorf("update target %d: %w", targetID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}

func (s Service) DeleteTargetByID(ctx context.Context, missionID, targetID int) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err = s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		target, err := s.targetsRepository.One(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("get target by id: %w", err)
		}

		if target.IsCompleted {
			return apperrors.TargetAlreadyCompleted(targetID).Wrap("can't delete")
		}

		err = s.targetsRepository.Delete(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("delete target %d: %w", targetID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}

func (s Service) AddTargetNote(ctx context.Context, missionID, targetID int, contents []string) (err error) {
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		mission, err := s.missionsRepository.One(ctx, missionID)
		if err != nil {
			return fmt.Errorf("get mission %d: %w", missionID, err)
		}

		if mission.IsCompleted {
			return apperrors.MissionAlreadyCompleted(targetID)
		}

		target, err := s.targetsRepository.One(ctx, missionID, targetID)
		if err != nil {
			return fmt.Errorf("get target by id: %w", err)
		}

		if target.IsCompleted {
			return apperrors.TargetAlreadyCompleted(targetID)
		}

		err = s.notesRepository.Create(ctx, missionID, targetID, contents)
		if err != nil {
			return fmt.Errorf("add notes by target id %d: %w", targetID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("within transaction: %w", err)
	}

	return nil
}
