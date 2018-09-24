package main

import (
	"path"
	"testing"
)

func TestMatches(t *testing.T) {
	t.Run("non-error matches", func(t *testing.T) {
		targetDir := "/home/moot"
		cases := []struct {
			path    string
			pattern string
			match   bool
		}{
			{"/home/moot/public_html/indx.php", "public_html/*", true},
			{"/home/moot/public_html/indx.php", "www/*", false},
			{"/home/moot/.local/bin/go", "*/bin/*", true},
		}

		for _, tt := range cases {
			pattern := path.Join(targetDir, tt.pattern)
			m, err := match(pattern, tt.path)

			if err != nil {
				t.Logf("bad pattern: %s", pattern)
				t.Fail()
			}

			if m != tt.match {
				t.Logf("expected match: %s=%s>%t", pattern, tt.path, m)
				t.Fail()
			}
		}
	})
	t.Run("error matches", func(t *testing.T) {
		targetDir := "/home/moot"
		cases := []struct {
			path    string
			pattern string
		}{
			// character classes must be non-empty. currently the only example of a malformed match pattern I can find
			{"/home/moot/public_html/indx.php", "public_html/[]"},
		}

		for _, tt := range cases {
			pattern := path.Join(targetDir, tt.pattern)
			_, err := match(pattern, tt.path)

			if err == nil {
				t.Logf("expected bad pattern: %s", pattern)
				t.Fail()
			}
		}
	})
}

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
			r, err := parseRule(tt.in)

			if err != nil {
				t.Logf("parseRule threw error: %s", err)
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
			_, err := parseRule(tt.in)

			if err != tt.out {
				t.Logf("expected error: %s", tt.out)
				t.Logf("received error:: %s", err)
				t.Fail()
			}

		}
	})
}
