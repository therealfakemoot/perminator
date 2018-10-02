package perminator

import (
	"testing"
)

func TestParseRule(t *testing.T) {
	t.Run("valid rules", func(t *testing.T) {
		cases := []struct {
			in  string
			out Rule
		}{
			{"*/bin f0655", Rule{Pattern: "*/bin", Type: "f", Mode: 0655}},
			{"*public_html/* d0644", Rule{Pattern: "*public_html/*", Type: "d", Mode: 0644}},
		}

		for _, tt := range cases {
			r, err := ParseRule(tt.in)

			if err != nil {
				t.Logf("ParseRule threw error: %s", err)
				t.Fail()
			}

			if r.Type != tt.out.Type || r.Mode != tt.out.Mode || r.Pattern != tt.out.Pattern {
				t.Logf("Failing Rule: %+v", r)
				t.Logf("Expected Rule: %+v", tt.out)
				t.Fail()
			}
		}
	})

	t.Run("invalid rules", func(t *testing.T) {
		cases := []struct {
			in  string
			out error
		}{
			{"*/bin x0655", ErrBadFileType},
			{"*public_html/* d0999", ErrBadPerms},
		}

		for _, tt := range cases {
			_, err := ParseRule(tt.in)

			if err != tt.out {
				t.Logf("expected error: %s", tt.out)
				t.Logf("received error:: %s", err)
				t.Fail()
			}

		}
	})
}
