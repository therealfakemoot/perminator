package test

import (
	"os"
	"testing"

	p "github.com/therealfakemoot/perminator/src"
)

func BenchmarkMatch(b *testing.B) {
	dir := "/home/user/really/long/path/with/lots/of/dirs"
	r := p.Rule{
		Pattern: "/home/user/really/long/path/*/dirs",
		Type:    "d",
		Mode:    os.FileMode(0777),
	}

	fi := FileInfoMock{FName: "", Dir: true}
	b.Run("simple match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p.Match(r, dir, fi)
		}
	})
	b.Run("no match", func(b *testing.B) {
		r := p.Rule{
			Pattern: "/home/otheruser/really/long/path/*/dirs",
			Type:    "d",
			Mode:    os.FileMode(0777),
		}

		fi := FileInfoMock{FName: "", Dir: true}
		for i := 0; i < b.N; i++ {
			p.Match(r, dir, fi)
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
