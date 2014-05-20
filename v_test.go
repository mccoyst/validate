package validate

import (
	"fmt"
	"strconv"
	"testing"
)

func ExampleVWithTwoParams_Validate() {
	type X struct {
		A string `validate:"between[5|12]"`
		B string `validate:"between[4|7]"`
		C string
		D string
	}

	vd := make(V)

	vd["between"] = func(i interface{}, p []string) error {
		s := i.(string)

		switch len(p) {
		case 1:
			l, err := strconv.Atoi(p[0])

			if err != nil {
				panic(err) //"Param needs to be an int")
			}

			if len(s) != l {
				return fmt.Errorf("%q needs to be %d in length", s, l)
			}
		case 2:
			b, err := strconv.Atoi(p[0])
			if err != nil {
				panic(err) //"Param needs to be an int")
			}
			e, err := strconv.Atoi(p[1])
			if err != nil {
				panic(err) //"Param needs to be an int")
			}
			if len(s) < b || len(s) > e {
				return fmt.Errorf("%q needs to be between %d and %d in length", s, b, e)
			}
		}

		return nil
	}

	fmt.Println(vd.Validate(X{
		A: "hello there",
		B: "hi",
		C: "help me",
		D: "I am not validated",
	}))

	// Output: [field B is invalid: "hi" needs to be between 4 and 7 in length]
}

func ExampleVWithOneParam_Validate() {
	type X struct {
		A string
		B string
		C string `validate:"max_length[10]"`
		D string `validate:"max_length[10]"`
	}

	vd := make(V)

	vd["max_length"] = func(i interface{}, p []string) error {
		s := i.(string)

		l, err := strconv.Atoi(p[0])

		if err != nil {
			panic(err) //"Param needs to be an int")
		}

		if len(s) > l {
			return fmt.Errorf("%q needs to be less than %d in length", s, l)
		}

		return nil
	}

	fmt.Println(vd.Validate(X{
		A: "hello there",
		B: "hi",
		C: "help me",
		D: "I am not validated",
	}))

	// Output: [field D is invalid: "I am not validated" needs to be less than 10 in length]
}

func ExampleV_Validate() {
	type X struct {
		A string `validate:"long"`
		B string `validate:"short"`
		C string `validate:"long,short"`
		D string
	}

	vd := make(V)
	vd["long"] = func(i interface{}, p []string) error {
		s := i.(string)
		if len(s) < 5 {
			return fmt.Errorf("%q is too short", s)
		}
		return nil
	}
	vd["short"] = func(i interface{}, p []string) error {
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

func TestV_Validate_allgood(t *testing.T) {
	type X struct {
		A int `validate:"odd"`
	}

	vd := make(V)
	vd["odd"] = func(i interface{}, p []string) error {
		n := i.(int)
		if n&1 == 0 {
			return fmt.Errorf("%d is not odd", n)
		}
		return nil
	}

	errs := vd.Validate(X{
		A: 1,
	})

	if errs != nil {
		t.Fatalf("unexpected errors for a valid struct: %v", errs)
	}
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
	if errs[0].Error() != `field A is invalid: undefined validator: "oops"` {
		t.Fatal("wrong message for an undefined validator:", errs[0].Error())
	}
}

func TestV_Validate_multi(t *testing.T) {
	type X struct {
		A int `validate:"nonzero,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}, p []string) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}, p []string) error {
		n := i.(int)
		if n&1 == 0 {
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

func ExampleV_Validate_struct() {
	type X struct {
		A int `validate:"nonzero"`
	}

	type Y struct {
		X `validate:"struct,odd"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}, p []string) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}
	vd["odd"] = func(i interface{}, p []string) error {
		x := i.(X)
		if x.A&1 == 0 {
			return fmt.Errorf("%d is not odd", x.A)
		}
		return nil
	}

	errs := vd.Validate(Y{X{
		A: 0,
	}})

	for _, err := range errs {
		fmt.Println(err)
	}

	// Output: field A is invalid: should be nonzero
	// field X is invalid: 0 is not odd
}

func TestV_Validate_uninterfaceable(t *testing.T) {
	type X struct {
		a int `validate:"nonzero"`
	}

	vd := make(V)
	vd["nonzero"] = func(i interface{}, p []string) error {
		n := i.(int)
		if n == 0 {
			return fmt.Errorf("should be nonzero")
		}
		return nil
	}

	errs := vd.Validate(X{
		a: 0,
	})

	if len(errs) != 0 {
		t.Fatal("wrong number of errors for two failures:", errs)
	}
}

func TestV_Validate_nonstruct(t *testing.T) {
	vd := make(V)
	vd["wrong"] = func(i interface{}, p []string) error {
		return fmt.Errorf("WRONG: %v", i)
	}

	errs := vd.Validate(7)
	if errs != nil {
		t.Fatal("non-structs should always pass validation: %v", errs)
	}
}
