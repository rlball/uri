package uri

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

func isAlias(v reflect.Value) bool {
	if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr {
		return false
	}
	return strings.Contains(v.Type().String(), ".")
}

func implementsUnmarshaler(v reflect.Value) bool {
	return v.Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem())
}

func implementsMarshaler(v reflect.Value) bool {
	return v.Type().Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem())
}

func tryMarshal(v reflect.Value) (string, error) {
	// does it implement TextMarshaler?
	if implementsMarshaler(v) {
		b, err := v.Interface().(encoding.TextMarshaler).MarshalText()
		return string(b), err
	} else if v.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
		return v.Interface().(fmt.Stringer).String(), nil
	}
	return "", nil
}

func isZero(v reflect.Value) bool {
	if !v.CanInterface() {
		return false
	}
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

// parseURITag splits the passed tag value by comma and returns the first value. Comma separated values
// is only checked when using the "json" tag as the name.
func parseURITag(tv string) string {
	if !usingJSONTag {
		return tv
	}

	return strings.Split(tv, ",")[0]
}
