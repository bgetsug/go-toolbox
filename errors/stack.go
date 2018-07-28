package errors

import (
	"fmt"

	"github.com/go-stack/stack"
)

type CallStack stack.CallStack

func (cs CallStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, c := range cs {
				fmt.Fprintf(s, "\n%+v (%n)", c, c)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []stack.Call(cs))
		default:
			fmt.Fprintf(s, "%v", []stack.Call(cs))
		}
	case 's':
		fmt.Fprintf(s, "%s", []stack.Call(cs))
	}
}
