package dto

type CreateCatParams struct {
	Name       string
	Breed      string
	Experience int
	Salary     int
}

type UpdateCatParams struct {
	CatID           int
	Name            *string
	ExperienceYears *int
	Salary          *int
}

type GetMissionsParams struct {
	CatID *int
}

type CreateMissionParams struct {
	Targets []CreateTargetParams
}

type UpdateMissionParams struct {
	MissionID     int
	AssignedCatID *int
	IsCompleted   *bool
}

type UpdateTargetParams struct {
	MissionID   int
	TargetID    int
	IsCompleted *bool
}

type CreateTargetParams struct {
	Name    string
	Country string
}
