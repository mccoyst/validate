// © 2013 Steve McCoy under the MIT license.

/*
Package validate provides a type for automatically validating the fields of structs.

Any fields tagged with the key "validate" will be validated via a user-defined list of functions.
For example:

	type X struct {
		A string `validate:"long"`
		B string `validate:"short"`
		C string `validate:"long,proper"`
		D string
	}

As shown, multiple validators can be named by separating each with a comma.

The validators are defined in a map like so:

	vd := make(validate.V)
	vd["long"] = func(i interface{}) error {
		…
	}
	vd["short"] = func(i interface{}) error {
		…
	}
	…

Validate passes the values of the tagged fields to these functions,
which should return an error when they decide a value is invalid.

There is a reserved tag, "struct",
which can be used to automatically validate
the fields of a named or embedded struct field.
"struct" may be combined with user-defined validators.

Reflection is used to access the tags and fields,
so the usual caveats and limitations apply.
*/
package validate

import (
	"fmt"
	"reflect"
	"strings"
)

// V is a map of tag names to validators.
type V map[string]func(interface{}) error

// BadField is an error type containing a field name and associated error.
// This is the type returned from Validate.
type BadField struct {
	Field string
	Err   error
}

func (b BadField) Error() string {
	return fmt.Sprintf("field %s is invalid: %v", b.Field, b.Err)
}

// Validate accepts a struct (or a pointer) and returns a list of errors for all
// fields that are invalid. If all fields are valid, or s is not a struct type,
// Validate returns nil.
//
// Fields that are not tagged or cannot be interfaced via reflection
// are skipped.
func (v V) Validate(s interface{}) []error {
	return v.ValidateAndTag(s, "")
}

// ValidateAndTag behaves like Validate, but uses the value of the supplied
// nameTag instead of the field name when reporting errors. For example:
//
//	struct X{
//		Y int `json:"height" validate:"nonzero"`
//	}
//
//	errs := v.ValidateAndTag(X{}, "json")
//
// The returned BadField will contain "height" instead of "Y" in Field.
//
// When nameTag == "", ValidateAndTag behaves identically to Validate.
func (v V) ValidateAndTag(s interface{}, nameTag string) []error {
	return v.validateAndTagPrefix(s, nameTag, "")
}

func (v V) validateAndTagPrefix(s interface{}, nameTag string, prefix string) []error {
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	t := val.Type()
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	var errs []error

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := val.Field(i)
		if !fv.CanInterface() {
			continue
		}
		val := fv.Interface()
		tag := f.Tag.Get("validate")
		if tag == "" {
			continue
		}
		vts := strings.Split(tag, ",")

		for _, vt := range vts {
			name := f.Name
			if nameTag != "" {
				name = f.Tag.Get(nameTag)
			}

			if len(prefix) > 0 {
				name = prefix + "." + name
			}

			if vt == "struct" {
				errs2 := v.validateAndTagPrefix(val, nameTag, name)
				if len(errs2) > 0 {
					errs = append(errs, errs2...)
				}
				continue
			}

			vf := v[vt]
			if vf == nil {
				errs = append(errs, BadField{
					Field: name,
					Err:   fmt.Errorf("undefined validator: %q", vt),
				})
				continue
			}
			if err := vf(val); err != nil {
				errs = append(errs, BadField{name, err})
			}
		}
	}

	return errs
}
