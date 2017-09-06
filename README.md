Extension for github.com/urfave/cli

## Usage
```go
package clix

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/urfave/cli"
)

type genericType struct {
	s string
}

func (g *genericType) Set(value string) error {
	g.s = value
	return nil
}

func (g *genericType) String() string {
	return g.s
}

func init() {
	os.Setenv("GENERIC_ENV", "generic_env")
}

var config struct {
	GenericWithDefaultValue *genericType  `name:"generic" value:"generic_test" usage:"generic type with default value"`
	GenericWithEnvVar       *genericType  `name:"generic2" env:"GENERIC_ENV" usage:"generic type with environment variable"`
	GenericWithEnvVar2      *genericType  `name:"generic3,g" env:"GENERIC_ENV" usage:"generic type with environment variable override by -g"`
	BoolValue               bool          `name:"bool" usage:"flag with type bool"`
	BoolTValue              bool          `name:"bool_true" value:"true" usage:"flag with type bool that is true by default"`
	DurationValue           time.Duration `name:"duration" value:"1m30s" usage:"flag with type duration that is 1m30s by default"`
	Int64Value              int64         `name:"int64" value:"1000" usage:"flag with type int64 that is 1000 by default"`
	IntValue                int           `name:"int" value:"1000" usage:"flag with type int that is 1000 by default"`
	Uint64Value             uint64        `name:"uint64" value:"1000" usage:"flag with type uint64 that is 1000 by default"`
	UintValue               uint          `name:"uint" value:"1000" usage:"flag with type uint that is 1000 by default"`
	Int64SliceValue         []int64       `name:"int64_slice" value:"1000,200,3000" usage:"flag with type []int64"`
	IntSliceValue           []int         `name:"int_slice" value:"1000,200,3000" usage:"flag with type []int"`
	StringValue             string        `name:"string" value:"hello clix" usage:"flag with type string"`
	StringSliceValue        []string      `name:"string_slice" value:"hello, clix" usage:"flag with type []string"`
}

func main(t *testing.T) {
	app := cli.NewApp()
	app.Flags = clix.MakeFlags(&config)
	app.Before = clix.MakeParser(&config)
	app.Action = func(ctx *cli.Context) error {
		println(fmt.Sprintf("%+v", config))
		return nil
	}

	app.Run([]string{"mockcmd",
		"--generic", "2222222",
		"-g", "override_env",
		"--bool",
		"--duration", "2m40s",
		"--int64", "5000",
		"--string", "override",
		"--int_slice", "1", "2,", "3",
		"--string_slice", "clix", "is", "awesome",
	})
}
```

## LICENSE
[MIT LICENSE](./LICENSE)