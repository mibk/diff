# diff

Package diff implements methods for comparing objects and producing
edit scripts. The motivation to create the package was to be able
to use the diff output format in tests where the output of
[go-cmp](https://github.com/google/go-cmp) wasn't suitable. It
isn't optimized for performance and as of now, it is a non-goal.
[See the package documentation for more information](https://godoc.org/github.com/mibk/diff).

## Instalation

```
$ go get github.com/mibk/diff
```
