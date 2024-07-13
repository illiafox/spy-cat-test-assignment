package http

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/service/dto"
)

type CreateCatRequest struct {
	Name       string `json:"name"`
	Breed      string `json:"breed"`
	Experience int    `json:"experience"`
	Salary     int    `json:"salary"`
}

func (r CreateCatRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Experience, validation.Min(0), validation.Max(100)),
		validation.Field(&r.Breed, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Salary, validation.Required, validation.Min(1), validation.Max(10000000)),
	)
}

type AddTargetRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func (r AddTargetRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Country, validation.Required, is.CountryCode2),
	)
}

type CreateMissionRequest struct {
	Targets []AddTargetRequest `json:"targets"`
}

func (r CreateMissionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Targets, validation.Required, validation.Length(1, 3)),
	)
}

func (r CreateMissionRequest) Params() dto.CreateMissionParams {
	targets := make([]dto.CreateTargetParams, len(r.Targets))
	for i := range r.Targets {
		targets[i] = dto.CreateTargetParams{
			Name:    r.Targets[i].Name,
			Country: r.Targets[i].Country,
		}
	}

	return dto.CreateMissionParams{
		Targets: targets,
	}
}

type UpdateCatRequest struct {
	CatID           int     `json:"-"`
	Name            *string `json:"name"`
	ExperienceYears *int    `json:"experience"`
	Salary          *int    `json:"salary"`
}

func (r UpdateCatRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.NilOrNotEmpty, validation.Length(2, 100)),
		validation.Field(&r.ExperienceYears, validation.Min(2), validation.Max(100)),
		validation.Field(&r.Salary, validation.Min(1), validation.Max(10000000)),
	)
}

// Missions

type GetMissionsRequest struct {
	CatID *int `query:"cat_id"`
}

type UpdateMissionRequest struct {
	AssignedCatID *int `json:"assigned_cat_id"`
}

func (r UpdateMissionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AssignedCatID, validation.Min(0)),
	)
}

// Targets

type AddTargetsRequest struct {
	Targets []AddTargetRequest `json:"targets"`
}

func (r AddTargetsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Targets, validation.Length(1, 3)),
	)
}

func (r AddTargetsRequest) Params() []dto.CreateTargetParams {
	targets := make([]dto.CreateTargetParams, len(r.Targets))
	for i := range targets {
		targets[i] = dto.CreateTargetParams{
			Name:    r.Targets[i].Name,
			Country: r.Targets[i].Country,
		}
	}
	return targets
}

type AddTargetNotesRequest struct {
	Notes []string `json:"notes"`
}

func (r AddTargetNotesRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Notes,
			validation.Required,
			validation.Length(1, 10000),
			validation.Each(validation.Length(1, 10000)),
		),
	)
}
