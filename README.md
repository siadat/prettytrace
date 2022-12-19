# prettytrace

```go
prettytrace.Print()
```

or

```go
prettytrace.Fprint(os.Stdout)
```

That's it.

## Examples:

See the [examples](/examples/) directory for actual Go files.

### Example 1:

```go
package main

import "github.com/siadat/prettytrace"

func main() {
	prettytrace.Print()
}
```

Output:

```
1: running
      > main/main.go:6  main(...)
```

### Example 2:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/siadat/prettytrace"
)

func main() {
	defer func() {
		var p = recover()
		fmt.Println(p)
		prettytrace.Print()
	}()
	strings.Repeat(" ", -1)
}
```

Output:

```
strings: negative Repeat count
1: running
    > main/main.go:12        main.func1(...)
      strings/strings.go:538 Repeat(...)
    > main/main.go:14        main(...)
```
