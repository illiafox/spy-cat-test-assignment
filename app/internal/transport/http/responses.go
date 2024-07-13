package http

import (
	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
	"time"
)

type BaseResponse struct {
	Ok bool `json:"ok"`
}

type Cat struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	ExperienceYears int16     `json:"experience_years"`
	Breed           string    `json:"breed"`
	Salary          int       `json:"salary"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func CatFromModel(cat *models.Cat) Cat {
	return Cat(*cat)
}

type GetCatsResponse struct {
	BaseResponse
	Cats []Cat `json:"cats"`
}

type GetCatResponse struct {
	BaseResponse
	Cat Cat `json:"cat"`
}

type CreateCatResponse struct {
	BaseResponse
	ID int `json:"id"`
}

// Missions

type Mission struct {
	ID            int       `json:"id"`
	AssignedCatID *int      `json:"assigned_cat_id"`
	IsCompleted   bool      `json:"is_completed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func MissionFromModel(mission *models.Mission) Mission {
	return Mission(*mission)
}

type MissionFull struct {
	Mission
	Targets []Target `json:"targets"`
}

func MissionFullFromModel(missionFull *models.MissionFull) MissionFull {
	mission := MissionFromModel(missionFull.Mission)

	targets := make([]Target, len(missionFull.Targets))
	for i := range targets {
		targets[i] = TargetFromModel(missionFull.Targets[i])
	}

	return MissionFull{
		Mission: mission,
		Targets: targets,
	}
}

type GetMissionsResponse struct {
	BaseResponse
	Missions []Mission `json:"missions"`
}

type GetMissionResponse struct {
	BaseResponse
	Mission MissionFull `json:"mission"`
}

type CreateMissionResponse struct {
	BaseResponse
	ID int `json:"id"`
}

// Targets
type Target struct {
	ID          int       `json:"id"`
	MissionID   int       `json:"mission_id"`
	IsCompleted bool      `json:"is_completed"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TargetFromModel(target *models.Target) Target {
	return Target(*target)
}

type TargetFull struct {
	Target
	Notes []Note `json:"notes"`
}

func TargetFullFromModel(targetFull *models.TargetFull) TargetFull {
	target := TargetFromModel(targetFull.Target)

	notes := make([]Note, len(targetFull.Notes))
	for i := range notes {
		notes[i] = NoteFromModel(targetFull.Notes[i])
	}

	return TargetFull{
		Target: target,
		Notes:  notes,
	}
}

type GetMissionTargetsResponse struct {
	BaseResponse
	Targets []TargetFull `json:"targets"`
}

type GetTargetResponse struct {
	BaseResponse
	Target TargetFull `json:"target"`
}

// Notes

type Note struct {
	ID        int       `json:"-"`
	TargetID  int       `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NoteFromModel(note *models.Note) Note {
	return Note(*note)
}
