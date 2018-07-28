package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhones(t *testing.T) {
	a := assert.New(t)

	err := Validate.Var("+1234567890", PhoneTag)
	a.Nil(err)

	err = Validate.Var("(123) 456-7890", EmailTag)
	a.NotNil(err)
}
