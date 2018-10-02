package perminator

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	mock "github.com/therealfakemoot/perminator/test"
)

func TestMatches(t *testing.T) {
	fiDir := mock.FileInfoMock{FName: "", Dir: true}
	fiFile := mock.FileInfoMock{FName: "", Dir: false}
	basePath := "/home/user"

	t.Run("positive matches", func(t *testing.T) {
		cases := []struct {
			r      Rule
			target string
			fi     os.FileInfo
			match  bool
		}{
			{r: Rule{Pattern: "public_html/*.php", Type: "f", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/index.php"),
				fi:     fiFile,
				match:  true},
			{r: Rule{Pattern: "public_html/controllers/", Type: "d", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/controllers/"),
				fi:     fiDir,
				match:  true},
			{r: Rule{Pattern: "public_html/controllers/", Type: "a", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/controllers/"),
				fi:     fiDir,
				match:  true},
		}

		for _, tt := range cases {
			m, _ := Match(tt.r, tt.target, tt.fi)

			if m != tt.match {
				t.Logf("Expected match for path `%s`: %+v", tt.target, tt.r)
			}
		}
	})

	t.Run("negative matches", func(t *testing.T) {
		cases := []struct {
			r      Rule
			target string
			fi     os.FileInfo
			match  bool
		}{
			{r: Rule{Pattern: "public_html/controllers/", Type: "d", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/index.php"),
				fi:     fiFile,
				match:  true},
			{r: Rule{Pattern: "public_html/*.php", Type: "f", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/controllers/"),
				fi:     fiDir,
				match:  true},
			{r: Rule{Pattern: "public_html/assets/*", Type: "a", Mode: os.FileMode(0777)},
				target: filepath.Join(basePath, "public_html/controllers/"),
				fi:     fiDir,
				match:  true},
		}

		for _, tt := range cases {
			m, _ := Match(tt.r, tt.target, tt.fi)

			if m != tt.match {
				t.Logf("Expected match for path `%s`: %+v", tt.target, tt.r)
			}
		}
	})
	t.Run("glob syntax errors", func(t *testing.T) {
		fi := mock.FileInfoMock{FName: "", Dir: false}
		targetDir := "/home/user"
		cases := []struct {
			path    string
			pattern string
		}{
			// character classes must be non-empty. currently the only example of a malformed match pattern I can find
			{"/home/user/public_html/indx.php", "public_html/[]"},
		}

		for _, tt := range cases {
			pattern := path.Join(targetDir, tt.pattern)
			_, err := Match(Rule{
				Pattern: pattern,
				Type:    "a",
				Mode:    os.FileMode(0777),
			}, tt.path, fi)

			if err == nil {
				t.Logf("expected bad pattern: %s", pattern)
				t.Fail()
			}
		}
	})
}
