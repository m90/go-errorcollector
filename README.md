# go-errorcollector

[![Build Status](https://travis-ci.org/m90/go-errorcollector.svg?branch=master)](https://travis-ci.org/m90/go-errorcollector)
[![godoc](https://godoc.org/github.com/m90/go-errorcollector?status.svg)](http://godoc.org/github.com/m90/go-errorcollector)

> collect multiple errors in golang keeping the standard error interface

### Installation using go get

```sh
$ go get github.com/m90/go-errorcollector
```

### Usage

Instantiate a new collector using `New()`, collect errors using `Collect(error)` and compare against `nil` as usual:

```go
err := errorcollector.New()

for _, e := range elements {
    err := mutate(e)
    err.Collect(err) // nil will be skipped
}

if err != nil {
    // handle the error
}

```

Error messages will be concatenated:

```go
err := errorcollector.New()
err.Collect(errors.New("beep"))
err.Collect(errors.New("boop"))
msg := err.Error() // => "collected errors: beep, boop"
```

You can also collect another collector:

```go
err := errorcollector.New()
err.Collect(errors.New("rock"))
childErr := errorcollector.New()
childErr.Collect("n")
childErr.Collect("roll")
err.Collect(childErr)
msg := err.Error() // => "collected errors: rock, n, roll"
```

The collector satisfies the standard `error` interface:

```go
func checkForTypos(list []string) error {
    err := errorcollector.New()
    for _, string = range list {
        err.Collect(findMistakes(string))
    }
    return err
}
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)
