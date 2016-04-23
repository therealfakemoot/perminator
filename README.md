#Overview
massperms is written to be a performant, utilitarian application. Given a list of file glob patterns and a directory, apply file modes to matching files according to the provided glob patterns.

#Globs
The underlying library call uses the globbing syntax for [filepath.Match](https://golang.org/pkg/path/filepath/#Match).

#Configuration
The configuration file will follow this syntax:

```
*cache*/ 777
/home/*/etc/* 400
```

#Invocation
```
massperms
```
`massperms` without any arguments defaults to looking in the current user's home dir for a configuration file and uses the current working directory for applying the permissions patterns.

```
Usage of ./massperms:
  -config string
        Path to your massperms patterns file. (default "/home/username/.massperms.rc")
  -target string
        Path to directory to apply patterns to. (default "/home/username/massperms")
```

##TODO:

- [x] Configurable logging.
- [ ] Dry run for testing patterns.
- [ ] Profiling
- [ ] Tests?
