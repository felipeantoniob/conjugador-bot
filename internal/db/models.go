// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Gerund struct {
	Infinitive    string
	Gerund        string
	GerundEnglish sql.NullString
}

type Infinitive struct {
	Infinitive        string
	InfinitiveEnglish sql.NullString
}

type Mood struct {
	Mood        string
	MoodEnglish sql.NullString
}

type Pastparticiple struct {
	Infinitive            string
	Pastparticiple        string
	PastparticipleEnglish sql.NullString
}

type Tense struct {
	Tense        string
	TenseEnglish sql.NullString
}

type Verb struct {
	Infinitive  string
	Mood        string
	Tense       string
	VerbEnglish sql.NullString
	Form1s      sql.NullString
	Form2s      sql.NullString
	Form3s      sql.NullString
	Form1p      sql.NullString
	Form2p      sql.NullString
	Form3p      sql.NullString
}
