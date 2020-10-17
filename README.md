# envy

Environment variables utilities

## Usage

```go
package main

import (
	"flag"
	"fmt"
	"github.com/licaonfee/envy"
	"os"
)

func main() {
    //this value should be ignored, because -my-flag is set
    os.Setenv("MY_FLAG", "fooo")
    //this value is passed to my-time
    os.Setenv("MY_TIME", "1m")
    os.Args = []string{"my-bin", "-my-flag", "bar"}
    
    myflag := flag.String("my-flag", "", "my flag")
    mytime := flag.Duration("my-time", 0, "my-time")
    
    flag.Parse()
    envy.FillFlags(flag.CommandLine)

    fmt.Printf("my-flag: %s \n", *myflag)
    fmt.Printf("my-time: %v \n", *mytime)
}
```

