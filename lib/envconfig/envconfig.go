package envconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const ExpectStructTag = "envvar"

var (
	ErrConvertEnvvarFailed = errors.New("type conversion of environment variable failed")
)

// MergeEnvVars accepts a pointer to a struct and updates fields tagged with `env`
// in-place with values from environment variables. A struct field tagged `env:"BANANA"`
// would be populated with the environment variable BANANA.
func MergeEnvVars(structPtr interface{}, prefix string) (map[string]reflect.StructField, error) {
	if structPtr == nil {
		return nil, nil
	}

	replaceMap := make(map[string]reflect.StructField)

	// refuse to modify non-struct
	var isPtr, isStruct bool
	kind := reflect.ValueOf(structPtr).Kind()
	if kind == reflect.Ptr {
		isPtr = true
	}
	if reflect.Indirect(reflect.ValueOf(structPtr)).Kind() == reflect.Struct {
		isStruct = true
	}
	if !(isPtr && isStruct) {
		return nil, fmt.Errorf(
			"%w: expected structPtr to be pointer to struct, got %+v",
			ErrConvertEnvvarFailed, kind,
		)
	}

	ptr := reflect.ValueOf(structPtr)
	cfgElem := ptr.Elem()
	cfgType := cfgElem.Type()
	for i := 0; i < cfgElem.NumField(); i++ {
		fld := cfgElem.Field(i)
		fldType := cfgType.Field(i)

		if !fld.CanSet() || !fld.IsValid() {
			// either non-addressable (i.e. input provided wasn't a pointer)
			// or field isn't exported
			continue
		}

		tagged := fldType.Tag.Get(ExpectStructTag)
		if tagged == "" {
			// field has no "envvar" tag
			continue
		}

		key := prefix + tagged // case-sensitive!
		if found, ok := os.LookupEnv(key); ok {
			if err := assign(key, found, fld); err != nil {
				return nil, err
			} else {
				replaceMap[key] = fldType
			}
		}
	}
	return replaceMap, nil
}

func assign(key, val string, field reflect.Value) error {
	conversionErr := func(e error) error {
		return fmt.Errorf("%w (key=%s, err=%v)", ErrConvertEnvvarFailed, key, e)
	}

	kind := field.Kind()
	switch kind {

	case reflect.String:
		field.SetString(val)

	case reflect.Float64:
		if v, err := strconv.ParseFloat(val, 64); err != nil {
			return conversionErr(err)
		} else {
			field.SetFloat(v)
		}

	case reflect.Bool:
		// case-insensitive
		ok, err := strconv.ParseBool(val)
		if err != nil {
			return conversionErr(err)
		}
		field.SetBool(ok)

	case reflect.Int:
		if v, err := strconv.ParseInt(val, 0, 0); err != nil {
			return conversionErr(err)
		} else {
			field.SetInt(v)
		}

	default:
		err := fmt.Errorf("unsupported field type: %v", kind)
		return conversionErr(err)
	}
	return nil
}
