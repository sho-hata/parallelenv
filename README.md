# parallelenv

go linter that check whether environment variables are set in tests run in parallel

## Instruction

```sh
go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest
```

## Usage

```go
package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
  t.Parallel()
  t.Setenv("LANGUAGE", "Go")
}

```

```sh
go vet -vettool=(which parallelenv) ./...

./main_test.go:9:2: cannot set environment variables in parallel tests
./main_test.go:8:2: cannot set environment variables in parallel tests
```

### CircleCI

```yaml
- run:
    name: install parallelenv
    command: go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest

- run:
    name: run parallelenv
    command: go vet -vettool=`which parallelenv` ./...
```

### GitHub Actions

```yaml
- name: install parallelenv
  run: go install github.com/sho-hata/parallelenv/cmd/parallelenv@latest

- name: run parallelenv
  run: go vet -vettool=`which parallelenv` ./...
```
