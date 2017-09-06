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

func TestMakeFlags(t *testing.T) {
	expect := []cli.Flag{
		cli.GenericFlag{
			Name:   "generic",
			Hidden: false,
			Usage:  "generic type with default value",
			Value:  &genericType{"generic_test"},
		},
		cli.GenericFlag{
			Name:   "generic2",
			Hidden: false,
			Usage:  "generic type with environment variable",
			EnvVar: "GENERIC_ENV",
			Value:  &genericType{},
		},
		cli.GenericFlag{
			Name:   "generic3,g",
			Hidden: false,
			Usage:  "generic type with environment variable override by -g",
			EnvVar: "GENERIC_ENV",
			Value:  &genericType{},
		},
		cli.BoolFlag{
			Name:   "bool",
			Hidden: false,
			Usage:  "flag with type bool",
		},
		cli.BoolTFlag{
			Name:   "bool_true",
			Hidden: false,
			Usage:  "flag with type bool that is true by default",
		},
		cli.DurationFlag{
			Name:   "duration",
			Hidden: false,
			Usage:  "flag with type duration that is 1m30s by default",
			Value:  time.Minute + 30*time.Second,
		},
		cli.Int64Flag{
			Name:   "int64",
			Hidden: false,
			Usage:  "flag with type int64 that is 1000 by default",
			Value:  1000,
		},
		cli.IntFlag{
			Name:   "int",
			Hidden: false,
			Usage:  "flag with type int that is 1000 by default",
			Value:  1000,
		},
		cli.Uint64Flag{
			Name:   "uint64",
			Hidden: false,
			Usage:  "flag with type uint64 that is 1000 by default",
			Value:  1000,
		},
		cli.UintFlag{
			Name:   "uint",
			Hidden: false,
			Usage:  "flag with type uint that is 1000 by default",
			Value:  1000,
		},
		cli.Int64SliceFlag{
			Name:   "int64_slice",
			Hidden: false,
			Usage:  "flag with type []int64",
			Value:  &cli.Int64Slice{1000, 200, 3000},
		},
		cli.IntSliceFlag{
			Name:   "int_slice",
			Hidden: false,
			Usage:  "flag with type []int",
			Value:  &cli.IntSlice{1000, 200, 3000},
		},
		cli.StringFlag{
			Name:   "string",
			Hidden: false,
			Usage:  "flag with type string",
			Value:  "hello clix",
		},
		cli.StringSliceFlag{
			Name:   "string_slice",
			Hidden: false,
			Usage:  "flag with type []string",
			Value:  &cli.StringSlice{"hello", "clix"},
		},
	}
	flags := MakeFlags(&config)

	if len(flags) != len(expect) {
		t.Fatalf("difference flag length, expect=%d, got=%d", len(expect), len(flags))
	}

	for i := 0; i < len(flags); i++ {
		if g, e := flags[i], expect[i]; !reflect.DeepEqual(g, e) {
			t.Fatalf(`expect="%+v", got="%+v"`, e, g)
		}
	}

}

func TestMakeParser(t *testing.T) {
	app := cli.NewApp()
	app.Flags = MakeFlags(&config)
	app.Before = MakeParser(&config)
	app.Action = func(ctx *cli.Context) error {
		if got := config.GenericWithDefaultValue.String(); got != "2222222" {
			t.Fatalf("expect=%s, got=%s", "2222222", got)
		}

		if got := config.GenericWithEnvVar.String(); got != "generic_env" {
			t.Fatalf("expect=%s, got=%s", "generic_env", got)
		}

		if got := config.GenericWithEnvVar2.String(); got != "override_env" {
			t.Fatalf("expect=%s, got=%s", "override_env", got)
		}

		if got := config.BoolValue; got != true {
			t.Fatalf("expect=%t, got=%t", true, got)
		}

		if got := config.BoolTValue; got != true {
			t.Fatalf("expect=%t, got=%t", true, got)
		}

		if got := config.DurationValue; got != 2*time.Minute+40*time.Second {
			t.Fatalf("expect=%s, got=%s", time.Minute+30*time.Second, got)
		}

		if got := config.Int64Value; got != 5000 {
			t.Fatalf("expect=%d, got=%d", 5000, got)
		}

		if got := config.StringValue; got != "override" {
			t.Fatalf("expect=%d, got=%d", "override", got)
		}

		if got := config.IntSliceValue; reflect.DeepEqual([]int{1, 2, 3}, got) {
			t.Fatalf("expect=%+v, got=%+v", []int{1, 2, 3}, got)
		}

		if got := config.StringSliceValue; reflect.DeepEqual([]string{"clix", "is", "awesome"}, got) {
			t.Fatalf("expect=%+v, got=%+v", []string{"clix", "is", "awesome"}, got)
		}

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
