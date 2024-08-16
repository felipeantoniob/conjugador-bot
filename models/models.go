package models

type Verb struct {
	Infinitive  string
	Mood        string
	Tense       string
	VerbEnglish string
	Form1s      string
	Form2s      string
	Form3s      string
	Form1p      string
	Form2p      string
	Form3p      string
}

type TenseMood struct {
	Mood  string `json:"mood"`
	Tense string `json:"tense"`
}

type TenseMoodChoice struct {
	Name  string    `json:"name"`
	Value TenseMood `json:"value"`
}
