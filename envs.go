package envs

import (
	"golang.org/x/xerrors"
	"os"
	"reflect"
	"strconv"
)

// Load from environment variables through reading struct key.
// It compatible with Int, String, Bool only
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
