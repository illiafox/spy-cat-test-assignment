package http

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
)

type Handler struct {
	service Service
}

func (h Handler) GetCats(ctx *fiber.Ctx) error {
	cats, err := h.service.GetCats(ctx.Context())
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get cats: %w", err))
	}

	out := make([]Cat, len(cats))
	for i := range cats {
		out[i] = CatFromModel(cats[i])
	}

	var resp GetCatsResponse
	resp.Ok = true
	resp.Cats = out

	return ctx.JSON(resp)
}

func (h Handler) CreateCat(ctx *fiber.Ctx) error {
	var req CreateCatRequest
	if err := ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err := req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	catID, err := h.service.AddCat(ctx.Context(), dto.CreateCatParams(req))
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to add cat: %w", err))
	}

	var resp CreateCatResponse
	resp.Ok = true
	resp.ID = catID

	return ctx.Status(http.StatusCreated).JSON(resp)
}

func (h Handler) GetCatByID(ctx *fiber.Ctx) error {
	catID, err := ctx.ParamsInt("cat_id")
	if err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse cat id"))
	}

	cat, err := h.service.GetCatByID(ctx.Context(), catID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get cat: %w", err))
	}

	var resp GetCatResponse
	resp.Ok = true
	resp.Cat = CatFromModel(cat)

	return ctx.JSON(resp)
}

func (h Handler) UpdateCatByID(ctx *fiber.Ctx) error {
	catID, err := ctx.ParamsInt("cat_id")
	if err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse cat id"))
	}

	req := UpdateCatRequest{
		CatID: catID,
	}
	if err = ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err = req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	err = h.service.UpdateCatByID(ctx.Context(), dto.UpdateCatParams(req))
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to update cat: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

func (h Handler) DeleteCatByID(ctx *fiber.Ctx) error {
	catID, err := ctx.ParamsInt("cat_id")
	if err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse cat id"))
	}

	err = h.service.DeleteCatByID(ctx.Context(), catID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to delete cat: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

// Missions

func (h Handler) extractMissionID(ctx *fiber.Ctx) (int, error) {
	missionID, err := ctx.ParamsInt("mission_id")
	if err != nil {
		return -1, RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse mission id"))
	}

	return missionID, nil
}

func (h Handler) GetMissions(ctx *fiber.Ctx) error {
	var req GetMissionsRequest
	if err := ctx.QueryParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse query"))
	}

	missions, err := h.service.GetMissions(ctx.Context(), dto.GetMissionsParams(req))
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get missions: %w", err))
	}

	out := make([]Mission, len(missions))
	for i := range missions {
		out[i] = MissionFromModel(missions[i])
	}

	var resp GetMissionsResponse
	resp.Ok = true
	resp.Missions = out

	return ctx.JSON(resp)
}

func (h Handler) CreateMission(ctx *fiber.Ctx) error {
	var req CreateMissionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err := req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	missionID, err := h.service.CreateMission(ctx.Context(), req.Params())
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to create mission: %w", err))
	}

	var resp CreateMissionResponse
	resp.Ok = true
	resp.ID = missionID

	return ctx.Status(http.StatusCreated).JSON(resp)
}

func (h Handler) GetMissionByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	mission, err := h.service.GetMissionByID(ctx.Context(), missionID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get mission: %w", err))
	}

	var resp GetMissionResponse
	resp.Ok = true
	resp.Mission = MissionFullFromModel(mission)

	return ctx.JSON(resp)
}

func (h Handler) UpdateMissionByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	var req UpdateMissionRequest
	if err = ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err = req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	err = h.service.UpdateMissionByID(ctx.Context(), dto.UpdateMissionParams{
		MissionID:     missionID,
		AssignedCatID: req.AssignedCatID,
	})
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to update cat: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

func (h Handler) CompleteMissionByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	isCompleted := true

	err = h.service.UpdateMissionByID(ctx.Context(), dto.UpdateMissionParams{
		MissionID:   missionID,
		IsCompleted: &isCompleted,
	})
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to update cat: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

func (h Handler) DeleteMissionByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	err = h.service.DeleteMissionByID(ctx.Context(), missionID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to delete mission: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

// Targets

func (h Handler) extractTargetID(ctx *fiber.Ctx) (int, error) {
	targetID, err := ctx.ParamsInt("target_id")
	if err != nil {
		return -1, RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse target id"))
	}

	return targetID, nil
}

func (h Handler) GetMissionTargets(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	targets, err := h.service.GetTargetsByMissionID(ctx.Context(), missionID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to delete mission: %w", err))
	}

	out := make([]TargetFull, len(targets))
	for i := range targets {
		out[i] = TargetFullFromModel(targets[i])
	}

	var resp GetMissionTargetsResponse
	resp.Ok = true
	resp.Targets = out

	return ctx.JSON(resp)
}

func (h Handler) AddMissionTargets(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	var req AddTargetsRequest
	if err = ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err = req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	err = h.service.AddMissionTargets(ctx.Context(), missionID, req.Params())
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to create mission: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.Status(http.StatusCreated).JSON(resp)
}

func (h Handler) GetTargetByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	targetID, err := h.extractTargetID(ctx)
	if err != nil {
		return err
	}

	target, err := h.service.GetTargetByID(ctx.Context(), missionID, targetID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get target by id: %w", err))
	}

	var resp GetTargetResponse
	resp.Ok = true
	resp.Target = TargetFullFromModel(target)

	return ctx.JSON(resp)
}

func (h Handler) CompleteTargetByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	targetID, err := h.extractTargetID(ctx)
	if err != nil {
		return err
	}

	err = h.service.CompleteTargetByID(ctx.Context(), missionID, targetID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get target by id: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

func (h Handler) DeleteTargetByID(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	targetID, err := h.extractTargetID(ctx)
	if err != nil {
		return err
	}

	err = h.service.DeleteTargetByID(ctx.Context(), missionID, targetID)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get delete target: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}

// Notes

func (h Handler) AddTargetNote(ctx *fiber.Ctx) error {
	missionID, err := h.extractMissionID(ctx)
	if err != nil {
		return err
	}

	targetID, err := h.extractTargetID(ctx)
	if err != nil {
		return err
	}

	var req AddTargetNotesRequest
	if err = ctx.BodyParser(&req); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err).Wrap("parse body"))
	}

	if err = req.Validate(); err != nil {
		return RespondWithError(ctx, apperrors.InvalidRequest(err))
	}

	err = h.service.AddTargetNote(ctx.Context(), missionID, targetID, req.Notes)
	if err != nil {
		return RespondWithError(ctx, fmt.Errorf("failed to get delete target: %w", err))
	}

	var resp BaseResponse
	resp.Ok = true

	return ctx.JSON(resp)
}
