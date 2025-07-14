package core

type Action struct {
	ID          int64    `json:"id"`
	EventID     int64    `json:"eventId"`
	Reference   *string  `json:"reference,omitempty"`
	ReferenceId *int64   `json:"referenceId,omitempty"`
	Tags        []string `json:"tags"`
	Note        *string  `json:"note,omitempty"`
}

type ActionReference struct {
	Table *string `json:"table"`
	ID    *int64  `json:"id"`
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

type ActionResponse struct {
	ID        int64            `json:"id"`
	EventID   int64            `json:"eventId"`
	Reference *ActionReference `json:"reference"`
	Tags      []string         `json:"tags"`
	Note      *string          `json:"note,omitempty"`
}
