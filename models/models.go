package models

type TenseMood struct {
	Mood  string `json:"mood"`
	Tense string `json:"tense"`
}

type TenseMoodChoice struct {
	Name  string    `json:"name"`
	Value TenseMood `json:"value"`
}
