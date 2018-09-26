# Overview
perminator is written to be a performant, utilitarian application. Given a list of file glob patterns and a directory, apply file modes to matching files according to the provided glob patterns.


# Configuration
The configuration file will follow this syntax:

## Globs
The underlying library call uses the globbing syntax for [filepath.Match](https://golang.org/pkg/path/filepath/#Match).

```
*cache*/ f0777
/home/*/etc/* d1400
bar/foo* a0650
```

The first item shall be a glob pattern as above. The second item shall be an `f`, `a`, or `d` followed by an octal integer representing the filemode to set on each matched object. `f` matches files, `d` matches directories, and `a` matches all. Furthermore, rules are applied in the order they are loaded from the configuration. Higher priority or more specific rules should be placed closer to the bottom of the .perminatorrc.

Please note that all configuration directives are relative to the absolute version of the target path. For example, if `-targetDir = foo/` the absolute targetDir is `/path/to/the/targetDir/`. Given a rule `bar/*`, the resulting match pattern will be `/path/to/the/targetDir/bar/*`.

This can lead to unexpected behavior if your rule includes a given target directory. For example, a rule `bar/* d0655` and a `-targetDir = bar/` produces a match pattern of `/path/to/bar/bar/*`. If you wish for every target under `targetDir` to match, simply prefix the pattern with `*`: `* d0655` or `*/bin f0755`.

# Invocation
```
perminator
```
`perminator` without any arguments defaults to looking in the current user's home dir for a configuration file and uses the current working directory for applying the permissions patterns.

```
perminator --help
usage: perminator [<flags>]

Flags:
      --help   Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=/home/ndumas5/.perminator.rc
               Configuration file path.
  -d, --target=/home/ndumas5/work/perminator
               Target directory.
      --debug  Enable debugging output.
```
