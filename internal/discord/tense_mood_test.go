package discord

import (
	"testing"
)

func TestCreateTenseMoodMap(t *testing.T) {
	expected := map[string]TenseMood{
		"Present":                               {"Indicativo", "Presente"},
		"Preterite":                             {"Indicativo", "Pretérito"},
		"Imperfect":                             {"Indicativo", "Imperfecto"},
		"Conditional":                           {"Indicativo", "Condicional"},
		"Future":                                {"Indicativo", "Futuro"},
		"Present perfect":                       {"Indicativo", "Presente"},
		"Preterite perfect (Past anterior)":     {"Indicativo", "Pretérito anterior"},
		"Pluperfect (Past perfect)":             {"Indicativo", "Pluscuamperfecto"},
		"Conditional perfect":                   {"Indicativo", "Condicional perfecto"},
		"Future perfect":                        {"Indicativo", "Futuro perfecto"},
		"Present subjunctive":                   {"Subjuntivo", "Presente"},
		"Imperfect subjunctive":                 {"Subjuntivo", "Imperfecto"},
		"Future subjunctive":                    {"Subjuntivo", "Futuro"},
		"Present perfect subjunctive":           {"Subjuntivo", "Presente perfecto"},
		"Pluperfect (Past perfect) subjunctive": {"Subjuntivo", "Pluscuamperfecto"},
		"Future perfect subjunctive":            {"Subjuntivo", "Pretérito anterior"},
		"Imperative":                            {"Imperativo Afirmativo", "Presente"},
		"Negative Imperative":                   {"Imperativo Negativo", "Presente"},
	}

	result := createTenseMoodMap()

	if len(result) != len(expected) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(result))
		return
	}

	for name, value := range expected {
		if resultValue, ok := result[name]; !ok || resultValue != value {
			t.Errorf("For name %q, expected value %v, got %v", name, value, result[name])
		}
	}
}

func TestGetTenseMoodChoices(t *testing.T) {
	choices := getTenseMoodChoices()

	if len(choices) != len(tenseMoodChoices) {
		t.Errorf("Expected %d choices, got %d", len(tenseMoodChoices), len(choices))
		return
	}

	for i, choice := range choices {
		expectedName := tenseMoodChoices[i].Name
		if choice.Name != expectedName {
			t.Errorf("For index %d, expected name %q, got %q", i, expectedName, choice.Name)
		}
		if choice.Value != expectedName {
			t.Errorf("For index %d, expected value %q, got %q", i, expectedName, choice.Value)
		}
	}
}

func TestGetValueByName(t *testing.T) {
	tests := []struct {
		name     string
		expected TenseMood
		hasError bool
	}{
		{"Present", TenseMood{"Indicativo", "Presente"}, false},
		{"Imperfect", TenseMood{"Indicativo", "Imperfecto"}, false},
		{"Nonexistent", TenseMood{}, true},
	}

	for _, test := range tests {
		result, err := getValueByName(test.name)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected an error for name %q but got none", test.name)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect an error for name %q but got %v", test.name, err)
			}
			if result != test.expected {
				t.Errorf("For name %q, expected value %v, got %v", test.name, test.expected, result)
			}
		}
	}
}
