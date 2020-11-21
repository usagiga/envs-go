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
// - Passed `out` is NOT pointer or interface
// - Passed `out`'s element is NOT struct
// - Fields cannot be assignable value(readonly or unaddressable)
// - (int) Can't cast value
func Load(out interface{}) (err error) {
	outKind := reflect.TypeOf(out).Kind()
	if outKind != reflect.Ptr && outKind != reflect.Interface {
		return xerrors.New("Passed incompatible type. `out` must be pointer or interface.")
	}

	outType := reflect.TypeOf(out).Elem()
	outValue := reflect.ValueOf(out).Elem()
	outKind = outType.Kind()
	if outKind != reflect.Struct {
		return xerrors.New("Passed incompatible type. `out`'s element must be struct.")
	}

	for i := 0; i < outType.NumField(); i++ {
		fieldType := outType.Field(i)
		fieldKind := fieldType.Type.Kind()
		fieldVal := outValue.Field(i)

		if !fieldVal.CanSet() {
			return xerrors.New("Passed incompatible type. `out`'s fields must be assignable.")
		}

		// Get `envs` struct tag
		envKey := fieldType.Tag.Get("envs")

		// If "-", skip it
		if envKey == "-" {
			continue
		}

		// TODO : If pointer, process it's element

		// If nested, process it recursive
		if fieldKind == reflect.Struct {
			err = Load(fieldVal.Addr().Interface())
			if err != nil {
				return xerrors.Errorf("Error raised in nested value: %w", err)
			}
			continue
		}

		// If `envs` struct tag not set, skip it
		if envKey == "" {
			continue
		}

		// Look up env
		envVal, ok := os.LookupEnv(envKey)
		if !ok || envVal == "" {
			continue
		}

		// Set the value
		switch fieldKind {
		case reflect.String:
			newFieldVal := envVal
			fieldVal.SetString(newFieldVal)
		case reflect.Int:
			newFieldVal, err := strconv.Atoi(envVal)
			if err != nil {
				return xerrors.Errorf("Can't cast int value: %w", err)
			}

			fieldVal.SetInt(int64(newFieldVal))
		case reflect.Bool:
			newFieldVal := envVal == "true"
			fieldVal.SetBool(newFieldVal)
		}
	}

	return nil
}
