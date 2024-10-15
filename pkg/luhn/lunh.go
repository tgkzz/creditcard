package luhn

import "strconv"

type Luhn struct {
}

func NewLuhn() *Luhn {
	return &Luhn{}
}

func (l *Luhn) LuhnCheckCard(card string) (bool, error) {
	sum := 0
	alt := false

	for i := len(card) - 1; i >= 0; i-- {
		num, err := strconv.Atoi(string(card[i]))
		if err != nil {
			return false, err
		}
		if alt {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}
		sum += num
		alt = !alt
	}

	return sum%10 == 0, nil
}

func (l *Luhn) CalculateLuhnCheckDigit(card string) int {
	sum := 0
	alt := true

	for i := len(card) - 1; i >= 0; i-- {
		digit := int(card[i] - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}

	return (10 - (sum % 10)) % 10

}
