[![](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/myENA/debug)

# debug
Collection of GO packages designed to help with debugging your smelly code.

## Mutex

This debug [Mutex](https://godoc.org/github.com/myENA/debug/sync#Mutex) included in this package is designed to yell if
the lock is not released within a certain period of time.  The time and tick rates are configurable if you construct the
mutex using the [NewDebugMutex](https://godoc.org/github.com/myENA/debug/sync#NewDebugMutex) method.