package test

import (
	"path"
	"testing"

	p "github.com/therealfakemoot/perminator/src"
)

func BenchmarkMatch(b *testing.B) {
	dir := "/home/user/really/long/path/with/lots/of/dirs"
	pat := "/home/user/really/long/path/*/dirs"
	b.Run("simple match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p.Match(pat, dir)
		}
	})
	b.Run("no match", func(b *testing.B) {
		dir := "/home/user/really/long/path/with/lots/of/dirs"
		pat := "/home/otheruser/really/long/path/*/dirs"
		for i := 0; i < b.N; i++ {
			p.Match(pat, dir)
		}
	})
}

func BenchmarkParse(b *testing.B) {
	r := "bin/* f0655"
	b.Run("single rule", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = p.ParseRule(r)
		}
	})

	b.Run("parse error", func(b *testing.B) {
		r := "bin/* x0655"
		for i := 0; i < b.N; i++ {
			_, _ = p.ParseRule(r)
		}
	})
}

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
			m, err := p.Match(pattern, tt.path)

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
			_, err := p.Match(pattern, tt.path)

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
			out p.Rule
		}{
			{"*/bin f0655", p.Rule{Pattern: "*/bin", Type: "f", Mode: 0655}},
			{"*public_html/* d0644", p.Rule{Pattern: "*public_html/*", Type: "d", Mode: 0644}},
		}

		for _, tt := range cases {
			r, err := p.ParseRule(tt.in)

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
			{"*/bin x0655", p.ErrBadFileType},
			{"*public_html/* d0999", p.ErrBadPerms},
			// {"*public_html/* d0999", p.ErrBadPattern},
			// I currently have no idea how to force filepath.Match to throw a syntax related error.
			// Consider this TODO. Eventually.
		}

		for _, tt := range cases {
			_, err := p.ParseRule(tt.in)

			if err != tt.out {
				t.Logf("expected error: %s", tt.out)
				t.Logf("received error:: %s", err)
				t.Fail()
			}

		}
	})
}