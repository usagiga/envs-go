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
	kind := reflect.TypeOf(out).Kind()
	if kind != reflect.Ptr && kind != reflect.Interface {
		return xerrors.New("Passed incompatible type. `out` must be pointer or interface.")
	}

	elemT := reflect.TypeOf(out).Elem()
	elemV := reflect.ValueOf(out).Elem()

	elemK := elemT.Kind()
	if elemK != reflect.Struct {
		return xerrors.New("Passed incompatible type. `out`'s element must be struct.")
	}

	for i := 0; i < elemT.NumField(); i++ {
		fieldType := elemT.Field(i)
		fieldKind := fieldType.Type.Kind()
		fieldVal := elemV.Field(i)

		if !fieldVal.CanSet() {
			return xerrors.New("Passed incompatible type. `out`'s fields must be assignable.")
		}

		// Look up `env` struct tag
		envKey, ok := fieldType.Tag.Lookup("envs")
		if !ok || envKey == "" {
			continue
		}

		// If "-", skip it
		if envKey == "-" {
			continue
		}

		// Look up env
		envVal, ok := os.LookupEnv(envKey)
		if !ok || envKey == "" {
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
