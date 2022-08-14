# parallelenv

go linter that check whether environment variables are set in tests run in parallel

## Motivation
Calls to `t.Parallel()` and `t.Setenv()` in the test function will cause a panic, as environment variables cannot be set in tests that run in parallel.

We created this tool in the hope that it would be noticed by static analysis before the tests are run.

[![test_and_lint](https://github.com/sho-hata/parallelenv/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/sho-hata/parallelenv/actions/workflows/ci.yaml)
## installation

```sh
go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest
```

## Usage

### test file Example
```go
package main

import (
	"testing"
)

func TestMain(t *testing.T) {
  t.Parallel()
  t.Setenv("LANGUAGE", "Go")
}

```

### Analysis
```sh
$ parallelenv ./...

./main_test.go:7:2: cannot set environment variables in parallel tests
./main_test.go:8:2: cannot set environment variables in parallel tests
```

## CI
### CircleCI

```yaml
- run:
    name: install parallelenv
    command: go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest

- run:
    name: run parallelenv
    command: parallelenv ./...
```

### GitHub Actions

```yaml
- name: install parallelenv
  run: go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest

- name: run parallelenv
  run: parallelenv ./...
```
