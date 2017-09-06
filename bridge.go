package clix

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// Errors
var (
	ErrNotPointer = errors.New("clix: pointer required")
)

var durationType = reflect.TypeOf(time.Duration(0))

// MakeFlags build cli.Flag slice from a struct
func MakeFlags(v interface{}) []cli.Flag {
	c := reflect.TypeOf(v)
	for c.Kind() == reflect.Ptr {
		c = c.Elem()
	}

	flags := []cli.Flag{}
	count := c.NumField()
	for i := 0; i < count; i++ {
		field := c.Field(i)
		flag := makeFlag(field)
		if flag != nil {
			flags = append(flags, flag)
		}
	}

	return flags
}

// MakeParser returns a cli hook for parse command line arguments
func MakeParser(v interface{}) cli.BeforeFunc {
	hook := func(ctx *cli.Context) error {
		store := reflect.ValueOf(v)
		if store.Kind() != reflect.Ptr {
			return ErrNotPointer
		}
		for store.Kind() == reflect.Ptr {
			store = store.Elem()
		}

		storeType := store.Type()
		for storeType.Kind() == reflect.Ptr {
			storeType = storeType.Elem()
		}

		count := storeType.NumField()
		for i := 0; i < count; i++ {
			if err := fillField(ctx, storeType.Field(i), store.Field(i)); err != nil {
				return err
			}
		}

		return nil
	}

	return hook
}

func fillField(ctx *cli.Context, structField reflect.StructField, field reflect.Value) error {
	name := strings.TrimSpace(structField.Tag.Get("name"))
	if name == "" {
		return nil
	}

	parts := strings.Split(name, ",")
	if len(parts) < 1 {
		return fmt.Errorf("clix: invalid name tag `%s`", name)
	}

	name = strings.TrimSpace(parts[0])

	if field.Kind() == reflect.Ptr {
		v := ctx.Generic(name)
		field.Set(reflect.ValueOf(v))
		return nil
	}

	if !field.CanSet() {
		return fmt.Errorf("clix: field `%s` is not changable", structField.Name)
	}

	var (
		filedType = structField.Type
		v         interface{}
	)
	// duration flag
	if filedType == durationType {
		v = ctx.Duration(name)
	} else {
		switch filedType.Kind() {
		case reflect.String:
			v = ctx.String(name)

		case reflect.Bool:
			value := strings.TrimSpace(structField.Tag.Get("value"))
			var t = false
			if value != "" {
				if r, err := strconv.ParseBool(value); err != nil {
					panic(err)
				} else {
					t = r
				}
			}

			if t {
				v = ctx.Bool(name)
			} else {
				v = ctx.BoolT(name)
			}

		case reflect.Float64:
			v = ctx.Float64(name)

		case reflect.Int64:
			v = ctx.Int64(name)

		case reflect.Int:
			v = ctx.Int(name)

		case reflect.Uint64:
			v = ctx.Uint64(name)

		case reflect.Uint:
			v = ctx.Uint(name)

		case reflect.Slice:
			switch filedType.Elem().Kind() {
			case reflect.Int:
				v = ctx.IntSlice(name)

			case reflect.Int64:
				v = ctx.Int64Slice(name)

			case reflect.String:
				v = ctx.StringSlice(name)

			default:
				return fmt.Errorf("unsupported slice type: %s", filedType.Elem().Kind().String())
			}

		default:
			return fmt.Errorf("unsupported field type: %s", structField.Name)
		}
	}

	field.Set(reflect.ValueOf(v))
	return nil
}

// makeFlag build single flag via field
func makeFlag(filed reflect.StructField) cli.Flag {
	tags := filed.Tag
	var (
		name   = strings.TrimSpace(tags.Get("name"))
		usage  = strings.TrimSpace(tags.Get("usage"))
		env    = strings.TrimSpace(tags.Get("env"))
		value  = strings.TrimSpace(tags.Get("value"))
		hidden = false
	)

	if name == "" {
		return nil
	}

	if h := strings.TrimSpace(tags.Get("hidden")); h != "" {
		if t, err := strconv.ParseBool(h); err != nil {
			panic(err)
		} else {
			hidden = t
		}
	}

	fieldType := filed.Type

	// duration flag
	if fieldType == durationType {
		flag := cli.DurationFlag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}

		if value != "" {
			duration, err := time.ParseDuration(value)
			if err != nil {
				panic(err)
			}
			flag.Value = duration
		}
		return flag
	}

	// cli flags
	var flag cli.Flag
	switch fieldType.Kind() {
	case reflect.String:
		flag = cli.StringFlag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
			Value:  value,
		}

	case reflect.Bool:
		var t = false
		if value != "" {
			r, err := strconv.ParseBool(value)
			if err != nil {
				panic(err)
			}
			t = r
		}
		if t {
			flag = cli.BoolTFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
			}
		} else {
			flag = cli.BoolFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
			}
		}

	case reflect.Float64:
		f := cli.Float64Flag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}
		if value != "" {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic(err)
			}
			f.Value = v
		}
		flag = f

	case reflect.Int64:
		f := cli.Int64Flag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}
		if value != "" {
			v, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				panic(err)
			}
			f.Value = v
		}
		flag = f

	case reflect.Int:
		f := cli.IntFlag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}
		if value != "" {
			v, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				panic(err)
			}
			f.Value = int(v)
		}
		flag = f

	case reflect.Uint64:
		f := cli.Uint64Flag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}
		if value != "" {
			v, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				panic(err)
			}
			f.Value = v
		}
		flag = f

	case reflect.Uint:
		f := cli.UintFlag{
			Name:   name,
			Usage:  usage,
			Hidden: hidden,
			EnvVar: env,
		}
		if value != "" {
			v, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				panic(err)
			}
			f.Value = uint(v)
		}
		flag = f

	case reflect.Slice:
		switch fieldType.Elem().Kind() {
		case reflect.Int:
			f := cli.IntSliceFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
			}
			if value != "" {
				f.Value = &cli.IntSlice{}
				parts := strings.Split(value, ",")
				for _, s := range parts {
					s = strings.TrimSpace(s)
					if err := f.Value.Set(s); err != nil {
						panic(err)
					}
				}
			}
			flag = f

		case reflect.Int64:
			f := cli.Int64SliceFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
			}
			if value != "" {
				f.Value = &cli.Int64Slice{}
				parts := strings.Split(value, ",")
				for _, s := range parts {
					s = strings.TrimSpace(s)
					if err := f.Value.Set(s); err != nil {
						panic(err)
					}
				}
			}
			flag = f

		case reflect.String:
			f := cli.StringSliceFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
			}
			if value != "" {
				f.Value = &cli.StringSlice{}
				parts := strings.Split(value, ",")
				for _, s := range parts {
					s = strings.TrimSpace(s)
					if err := f.Value.Set(s); err != nil {
						panic(err)
					}
				}
			}
			flag = f

		default:
			panic("unsupported slice type")
		}

	default:
		if fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct {
			genType := fieldType.Elem()
			if genType.Kind() == reflect.Ptr {
				panic("clix: nest pointer type does not supported so far: " + fieldType.String())
			}

			instance := reflect.New(genType)
			v, ok := instance.Interface().(cli.Generic)
			if !ok {
				panic("clix: field type can not be converted to cli.Generic")
			}
			v.Set(value)
			flag = cli.GenericFlag{
				Name:   name,
				Usage:  usage,
				Hidden: hidden,
				EnvVar: env,
				Value:  v,
			}
		} else {
			panic("unsupported field type: " + fieldType.String())
		}
	}

	return flag
}
