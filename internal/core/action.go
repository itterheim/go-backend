package core

type Action struct {
	ID          int64    `json:"id"`
	EventID     int64    `json:"eventId"`
	Reference   *string  `json:"reference"`
	ReferenceId *string  `json:"referenceId"`
	Tags        []string `json:"tags"`
	Note        *string  `json:"note"`
}

type ActionReference struct {
	Table string `json:"table"`
	ID    int64  `json:"id"`
}

type CreateActionRequest struct {
	EventID   int64            `json:"eventId"`
	Reference *ActionReference `json:"reference"`
	Tags      []string         `json:"tags"`
	Note      *string          `json:"note"`
}

type UpdateActionRequest struct {
	ID        int64            `json:"id"`
	EventID   int64            `json:"eventId"`
	Reference *ActionReference `json:"reference"`
	Tags      []string         `json:"tags"`
	Note      *string          `json:"note"`
}
