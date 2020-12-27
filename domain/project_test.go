package domain

import "testing"

func TestProject_Test_Validation(t *testing.T) {
	testCases := map[string]struct {
		id string
		name string
	}{
		"empty ID": {
			name: "jfslfjal",
		},

		"empty name": {
			id: "jfslfjal",
		},
	}

	for _, tc := range testCases {
		_, err := NewProject(tc.id, tc.name)
		if err == nil {
			t.Error("expected that the validation fails but got no error")
		}
	}
}

