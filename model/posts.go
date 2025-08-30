package model

type Posts struct {
	ID      int    `json:"id"`
	Title   string `json:"title" validate:"required, min=5, max=200"`
	Content string `json:"content" validate:"required, min=50"`
	Date    string `json:"date" validate:"required"`
	Author  string `json:"author,omitempty"`
}
