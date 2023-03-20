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

## With mapstructure

```go
package main

import (
	"fmt"
	"log"

	"github.com/licaonfee/envy"
	"github.com/mitchellh/mapstructure"
)

type config struct {
	Path  string   `mapstructure:"path"`
	Users []string `mapstructure:"users"`
}

// Environment
// APP_PATH=log.log
// APP_USERS_0=one_user
// APP_USERS_1=other_user

func main() {
	var cfg config
	m := envy.FillMap(envy.NewOsEnv(), envy.FilterPrefix("APP_"))
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result: &cfg,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := dec.Decode(m); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", cfg)
}

```
