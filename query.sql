
-- name: GetVerbByInfinitiveMoodTense :one
SELECT * FROM verbs
WHERE infinitive = ? AND mood = ? AND tense = ?;
