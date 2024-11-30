-- name: GetVerbByInfinitiveMoodTense :one
SELECT
    infinitive,
    mood,
    tense,
    verb_english,
    form_1s,
    form_2s,
    form_3s,
    form_1p,
    form_2p,
    form_3p
FROM verbs
WHERE infinitive = ? AND mood = ? AND tense = ?;
