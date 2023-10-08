package handler

type TodoPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUpdatePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
