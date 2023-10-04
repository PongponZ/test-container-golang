package entity

type Todo struct {
	ID          int    `json:"id" bson:"_id,omitempty" `
	Title       string `json:"title" bson:"title,omitempty" `
	Description string `json:"description" bson:"description,omitempty" `
}
