package models

type (
	FileData struct {
		Location   string     `json:"location"`
		Dictionary []KeyValue `json:"dictionary"`
	}

	KeyValue struct {
		Key string      `json:"key"`
		Loc KeyValueLoc `json:"loc"`
	}

	KeyValueLoc struct {
		En string `json:"en"`
		Ru string `json:"ru"`
	}
)
