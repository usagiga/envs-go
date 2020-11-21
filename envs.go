package envs

import (
	"golang.org/x/xerrors"
	"os"
	"reflect"
	"strconv"
)

// Load from environment variables specified by "envs" struct tag.
// It compatible with int, string, bool.
//
// Some fields are ignored:
// - Fields have Incompatible types
// - Fields have `envs:"-"` struct tag
// - Fields don't have `envs` struct tag
// - Specified environment variables has no value
//
// Some values are passed, it raises error:
// - (int) Can't cast value
//
// Some values are passed, it raises panic:
// - Passed `out` is NOT pointer or interface
// - Passed `out`'s element is NOT struct
// - Fields cannot be assignable value(readonly or unaddressable)
func Load(out interface{}) (err error) {
	t := reflect.TypeOf(out).Elem()
	v := reflect.ValueOf(out).Elem()
	fields := t.NumField()

	for i := 0; i < fields; i++ {
		fType := t.Field(i)
		fTypeKind := fType.Type.Kind()
		fVal := v.Field(i)

		// Look up `env` struct tag
		envKey, ok := fType.Tag.Lookup("envs")
		if !ok || envKey == "" {
			continue
		}

		// If "-", skip it
		if envKey == "-" {
			continue
		}

		// Look up env
		envStrVal, ok := os.LookupEnv(envKey)
		if !ok || envKey == "" {
			continue
		}

		// Set the value
		switch fTypeKind {
		case reflect.String:
			envRefVal := reflect.ValueOf(envStrVal)
			fVal.Set(envRefVal)
		case reflect.Int:
			// Cast
			envIntVal, err := strconv.Atoi(envStrVal)
			if err != nil {
				return xerrors.Errorf("Can't cast int value: %w", err)
			}
			envRefVal := reflect.ValueOf(envIntVal)

			fVal.Set(envRefVal)
		case reflect.Bool:
			envBoolVal := envStrVal == "true"
			envRefVal := reflect.ValueOf(envBoolVal)

			fVal.Set(envRefVal)
		}
	}

	return nil
}
