package perminator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	Pattern string
	Type    string
	Mode    os.FileMode
}

type RuleSet []Rule

func LoadRules(path string) (RuleSet, error) {
	var rs RuleSet

	conf, err := ioutil.ReadFile(path)

	if err != nil {
		log.Print("error loading config file")
		return rs, err
	}

	lines := strings.Split(string(conf), "\n")

	log.Printf("raw config: %q", string(conf))
	log.Printf("config lines: %+v", lines)
	for _, line := range strings.Split(string(conf), "\n") {
		if line == "" {
			continue
		}
		r, err := ParseRule(line)
		if err != nil {
			log.Printf("bad rule encountered, %q: %s", line, err)
			continue
		}
		rs = append(rs, r)
	}

	return rs, nil
}

func ParseRule(s string) (Rule, error) {
	var (
		pattern string
		fstype  string
		rawMode string
		r       Rule
	)

	_, err := fmt.Sscanf(s, "%s %1s%s", &pattern, &fstype, &rawMode)
	if err != nil {
		return r, err
	}

	switch fstype {
	case "d", "f", "a":
	default:
		return r, ErrBadFileType

	}

	i, err := strconv.ParseUint(rawMode, 8, 32)

	if err != nil {
		return r, ErrBadPerms
	}

	r = Rule{
		Pattern: pattern,
		Type:    fstype,
		Mode:    os.FileMode(uint32(i)),
	}

	return r, nil
}
