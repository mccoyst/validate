package validate

import "fmt"

func ExampleValidate() {
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
