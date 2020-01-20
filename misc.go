package termimg

import (
	"fmt"
	"reflect"
	"regexp"
)

func mustScanSubexps(ptn *regexp.Regexp, into ...interface{}) {
	if err := scanSubexps(ptn, into...); err != nil {
		panic(err)
	}
}

func scanSubexps(ptn *regexp.Regexp, into ...interface{}) error {
	names := ptn.SubexpNames()

	if len(into)%2 != 0 {
		return fmt.Errorf("scanSubexps must contain alternating 'string', '*int' args")
	}

	for i := 0; i < len(into); i += 2 {
		intoName := reflect.ValueOf(into[i])
		intoVal := reflect.ValueOf(into[i+1])

		if intoName.Kind() != reflect.String {
			return fmt.Errorf("scanSubexps must contain alternating 'string', '*int' args; found %q at index %d", intoName.Kind(), i)
		}
		if intoVal.Kind() != reflect.Ptr {
			return fmt.Errorf("scanSubexps must contain alternating 'string', '*int' args; found %q at index %d", intoVal.Kind(), i+1)
		}
		if intoVal.Elem().Kind() != reflect.Int {
			return fmt.Errorf("scanSubexps must contain alternating 'string', '*int' args; found %q at index %d", intoVal.Elem().Kind(), i+1)
		}

		intoNameStr := intoName.Interface().(string)
		intoValPtr := intoVal.Interface().(*int)
		var found bool
		for idx, n := range names {
			if n == intoNameStr {
				*intoValPtr = idx
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("scanSubexps could not find subexp %q", intoNameStr)
		}
	}

	return nil
}
