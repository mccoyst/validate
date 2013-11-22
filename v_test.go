package validate

import (
	"fmt"
	"testing"
)

func ExampleV_Validate() {
	type X struct {
		A string `validate:"long"`
		B string `validate:"short"`
		C string `validate:"long,short"`
		D string
	}

	vd := make(V)
	vd["long"] = func(i interface{}) error {
		s := i.(string)
		if len(s) < 5 {
			return fmt.Errorf("%q is too short", s)
		}
		return nil
	}
	vd["short"] = func(i interface{}) error {
		s := i.(string)
		if len(s) >= 5 {
			return fmt.Errorf("%q is too long", s)
		}
		return nil
	}

	fmt.Println(vd.Validate(X{
		A: "hello there",
		B: "hi",
		C: "help me",
		D: "I am not validated",
	}))

	// Output: [field C is invalid: "help me" is too long]
}

func TestV_Validate_undef(t *testing.T) {
	type X struct {
		A string `validate:"oops"`
	}

	vd := make(V)

	errs := vd.Validate(X{
		A: "oh my",
	})

	if len(errs) == 0 {
		t.Fatal("no errors returned for an undefined validator")
	}
	if len(errs) > 1 {
		t.Fatalf("too many errors returns for an undefined validator: %v", errs)
	}
	if errs[0].Error() != `field A has an undefined validator: "oops"` {
		t.Fatal("wrong message for an undefined validator:", errs[0].Error())
	}
}

func TestV_Validate_multi(t *testing.T) {
	type X struct {
		A int `validate:"nonzero,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}) error {
		n := i.(int)
		if n & 1 == 0 {
			return fmt.Errorf("%d is not odd", n)
		}
		return nil
	}
	errs := vd.Validate(X{
		A: 0,
	})

	if len(errs) != 2 {
		t.Fatal("wrong number of errors for two failures: %v", errs)
	}
	if errs[0].Error() != "field A is invalid: should be nonzero" {
		t.Fatal("first error should be nonzero:", errs[0])
	}
	if errs[1].Error() != "field A is invalid: 0 is not odd" {
		t.Fatal("second error should be odd:", errs[1])
	}
}
