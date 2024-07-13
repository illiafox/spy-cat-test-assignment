package models

import (
	"time"
)

type Cat struct {
	ID              int       `db:"id"`
	Name            string    `db:"name"`
	ExperienceYears int16     `db:"experience_years"`
	Breed           string    `db:"breed"`
	Salary          int       `db:"salary"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type Mission struct {
	ID            int       `db:"id"`
	AssignedCatID *int      `db:"assigned_cat_id"`
	IsCompleted   bool      `db:"is_completed"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type MissionFull struct {
	*Mission
	Targets []*Target
}

type Target struct {
	ID          int       `db:"id"`
	MissionID   int       `db:"mission_id"`
	IsCompleted bool      `db:"is_completed"`
	Name        string    `db:"name"`
	Country     string    `db:"country"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type TargetFull struct {
	*Target
	Notes []*Note
}

type Note struct {
	ID        int       `db:"id"`
	TargetID  int       `db:"target_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
