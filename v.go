// © 2013 Steve McCoy under the MIT license.

/*
Package validate provides a type for automatically validating the fields of structs.

Any fields tagged with the key "validate" will be validated via a user-defined list of validators.
For example:

	type X struct {
		A string `validate:"long"
		B string `validate:"short"
		C string `validate:"long,proper"
		D string
	}

The validators are defined in a map like so:

	vd := make(validate.V)
	vd["long"] = func(i interface{}) error {
		…
	}
	vd["short"] = func(i interface{}) error {
		…
	}
	…

Reflection is used to access the tags and fields, so the usual caveats and limitations apply.
*/
package validate

import (
	"fmt"
	"reflect"
	"strings"
)

// V is a map of tag names to validators.
type V map[string]func(interface{}) error

// Validate accepts a struct and returns a list of errors for all
// fields that are invalid. If all fields are valid, or s is not a struct type,
// Validate returns nil.
//
// Fields that are not tagged or cannot be interfaced via reflection
// are skipped.
func (v V) Validate(s interface{}) []error {
	t := reflect.TypeOf(s)
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	val := reflect.ValueOf(s)
	errs := []error{}

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
			vf := v[vt]
			if vf == nil {
				errs = append(errs, fmt.Errorf("field %s has an undefined validator: %q", f.Name, vt))
				continue
			}
			if err := vf(val); err != nil {
				errs = append(errs, fmt.Errorf("field %s is invalid: %v", f.Name, err))
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
