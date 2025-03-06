package models

import "encoding/json"

type Recipe struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	ImageURL    string          `json:"image_url"`
	CookingTime int             `json:"cooking_time"`
	Ingredients json.RawMessage `json:"ingredients"`
	Steps       json.RawMessage `json:"steps"`
}
