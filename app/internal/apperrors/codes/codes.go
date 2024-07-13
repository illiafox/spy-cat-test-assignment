package codes

type Code string

const (
	Internal                  Code = "INTERNAL_ERROR"
	InvalidRequest            Code = "INVALID_REQUEST"
	CatNotFound               Code = "CAT_NOT_FOUND"
	MissionNotFound           Code = "MISSION_NOT_FOUND"
	TargetNotFound            Code = "TARGET_NOT_FOUND"
	MissionAlreadyCompleted   Code = "MISSION_ALREADY_COMPLETED"
	CatAlreadyAssigned        Code = "CAT_ALREADY_ASSIGNED"
	TargetAlreadyCompleted    Code = "TARGET_ALREADY_COMPLETED"
	AllTargetsAreNotCompleted Code = "ALL_TARGETS_ARE_NOT_COMPLETED"
)
