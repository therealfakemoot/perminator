#Overview
massperms is written to be a performant, utilitarian application. Given a list of file glob patterns and a directory, apply file modes to matching files according to the provided glob patterns.

#Globs
The underlying library call uses the globbing syntax for [filepath.Match](https://golang.org/pkg/path/filepath/#Match).

#Configuration
The configuration file will follow this syntax:

```
*cache*/ f777
/home/*/etc/* d400
bar/foo* a650
```

The first item shall be a glob pattern as above. The second item shall be an `f`, `a`, or `d` followed by a three digit integer representing the filemode to set on each matched object. `f` matches files, `d` matches directories, and `a` matches all. Furthermore, rules are applied in the order they are loaded from the configuration. Higher priority or more specific rules should be placed closer to the bottom of the .masspermsrc.

#Invocation
```
massperms
```
`massperms` without any arguments defaults to looking in the current user's home dir for a configuration file and uses the current working directory for applying the permissions patterns.

```
massperms --help
usage: massperms [<flags>]

Flags:
      --help   Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=/home/ndumas5/.massperms.rc  
               Configuration file path.
  -d, --target=/home/ndumas5/work/massperms  
               Target directory.
      --debug  Enable debugging output.
```

##TODO:

- [x] Configurable logging.
- [ ] Dry run for testing patterns.
- [ ] Profiling
- [ ] Tests?
