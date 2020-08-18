# envs-go

Easiest environment variables binder in golang


## Installation

```sh
$ go get github.com/usagiga/envs-go
```


## Example

```go
package main

import (
	"github.com/usagiga/envs-go"
	"log"
)

type Config struct {
	// These fields will be loaded
	BoolVal   bool   `envs:"BOOL_VAL"`
	StringVal string `envs:"STRING_VAL"`
	IntVal    int    `envs:"BOOL_VAL"`

	// These fields will be ignored
	IgnoredVal1 string `envs:"-"`
	IgnoredVal2 string
}

func main() {
	config := &Config{}
	err := envs.Load(config)
	if err != nil {
		log.Fatalf("Can't load config: %+v", err)
	}

	log.Println("StringVal: ", config.StringVal)
}
```

If there's no `envs` struct tag, no value in specified environment keys or field applied `envs:"-"`, envs-go will ignore it.


## Features

- Compatible with `xerrors`
- Auto type detection

### Supported types

- `int`
- `string`
- `bool`


## Dependencies

- Go (1.15 or higher)
- [golang.org/x/xerrors](https://pkg.go.dev/golang.org/x/xerrors)


## License

MIT
