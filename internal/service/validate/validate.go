package validate

import (
	"creditcard/pkg/luhn"
	"errors"
	"fmt"
	"unicode"
)

type Validator interface {
	ValidateCard(card string) error
}

type Validate struct {
	luhnPkg *luhn.Luhn
}

var (
	ErrInvalidLength     = errors.New("invalid card length")
	ErrInvalidCardFormat = errors.New("invalid card format")
	ErrInvalidCharacter  = errors.New("invalid character")
)

func NewValidator() *Validate {
	return &Validate{
		luhnPkg: luhn.NewLuhn(),
	}
}

func (v *Validate) ValidateCard(card string) error {
	if len(card) < 13 {
		return ErrInvalidLength
	}

	for _, r := range card {
		if !unicode.IsDigit(r) {
			return ErrInvalidCardFormat
		}
	}

	luhnCheck, err := v.luhnPkg.LuhnCheckCard(card)
	if err != nil {
		return ErrInvalidCardFormat
	}

	if !luhnCheck {
		fmt.Println("INCORRECT")
		return nil
	}

	fmt.Println("OK")
	return nil
}
