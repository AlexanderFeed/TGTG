package httpapi

import "testing"

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		want      string
		wantValid bool
	}{
		{name: "normalizes case and whitespace", input: "  User@Example.RU ", want: "user@example.ru", wantValid: true},
		{name: "rejects missing domain", input: "user@", want: "user@", wantValid: false},
		{name: "rejects display address", input: "User <user@example.ru>", want: "user <user@example.ru>", wantValid: false},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, valid := normalizeEmail(test.input)
			if got != test.want || valid != test.wantValid {
				t.Fatalf("normalizeEmail(%q) = (%q, %v), want (%q, %v)", test.input, got, valid, test.want, test.wantValid)
			}
		})
	}
}

func TestNewCode(t *testing.T) {
	t.Parallel()
	for range 100 {
		code, err := newCode()
		if err != nil {
			t.Fatalf("newCode returned error: %v", err)
		}
		if !codePattern.MatchString(code) {
			t.Fatalf("newCode returned invalid code %q", code)
		}
	}
}
