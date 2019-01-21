# errtag

errtag generates error tags that take hierarchical structure into account.

## Install

```bash
$ go get -u github.com/pi9min/errtag
```

## How to use

Call `errtag.ErrorTag()`

e.g.

```go
// github.com/foo/bar/cmd/main.go
package main

import (
	"fmt"
	"github.com/pi9min/errtag"
)

func main() {
	fmt.Println(errtag.ErrorTag())
}
```

The execution result is as follows.

```bash
$ go run ./cmd/main.go
main.main
```

## License

See [LICENSE](/LICENSE)
