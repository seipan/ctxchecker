# ctxchecker

This oss checks if the function argument contains a context.

## Example
```go
package a

import (
	"context"
	"net/http"
	"testing"
)

func f1(ctx context.Context, flag string) {

}

func f2(flag string) {

}

func f3(w http.ResponseWriter, r *http.Request) {
}

func Test_f3(t *testing.T) {

}

```
result:
```
analysistest.go:454: a/a.go:13:8: unexpected diagnostic: no ctx
```

## Install
```
go install github.com/seipan/ctxchecker/cmd/ctxchecker
```

## Usage
```
go vet -vettool=`which ctxchecker` pkgname
```
