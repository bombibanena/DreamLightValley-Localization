package types

import "ddv_loc/pkg/models"

type LocFile []models.FileData

type (
	LocFileUpdates struct {
		New     LocFile               `json:"new"`
		Changes LocFileUpdatesChanges `json:"changes"`
		Removed LocFile               `json:"removed"`
	}

	LocFileUpdatesChanges struct {
		New     LocFile `json:"new"`
		Changes LocFile `json:"changes"`
		Removed LocFile `json:"removed"`
	}
)

func (u *LocFileUpdates) Any() bool {
	return len(u.New) > 0 || len(u.Changes.New) > 0 || len(u.Changes.Changes) > 0 || len(u.Changes.Removed) > 0 || len(u.Removed) > 0
}
