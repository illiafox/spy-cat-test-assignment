package schema

import (
	"time"

	"github.com/illiafox/spy-cat-test-assignment/app/internal/models"
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

func (c Cat) ToModel() *models.Cat {
	cat := models.Cat(c)
	return &cat
}

type Mission struct {
	ID            int       `db:"id"`
	AssignedCatID *int      `db:"assigned_cat_id"`
	IsCompleted   bool      `db:"is_completed"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (m Mission) ToModel() *models.Mission {
	mission := models.Mission(m)
	return &mission
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

func (t Target) ToModel() *models.Target {
	target := models.Target(t)
	return &target
}

type Note struct {
	ID        int       `db:"mission_id"`
	TargetID  int       `db:"target_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func (n Note) ToModel() *models.Note {
	note := models.Note(n)
	return &note
}
