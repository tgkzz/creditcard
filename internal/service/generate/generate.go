package generate

import (
	"creditcard/pkg/luhn"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
)

type Generator interface {
	Generate(card string, pick bool) ([]string, error)
}

type Generate struct {
	luhnPkg *luhn.Luhn
}

func NewGenerator() Generator {
	return &Generate{
		luhnPkg: luhn.NewLuhn(),
	}
}

var (
	ErrInvalidCardFormat = errors.New("invalid card format")
	ErrInvalidAstNum     = errors.New("invalid asterisks num")
)

func (g *Generate) Generate(card string, pick bool) ([]string, error) {
	if err := g.validateCardFormat(card); err != nil {
		return nil, err
	}

	astNum := strings.Count(card, "*")
	if astNum > 4 {
		return nil, ErrInvalidAstNum
	}

	combinations := g.generateCombinations(card, astNum)

	sort.Strings(combinations)

	validCards := make([]string, 0)
	for _, combination := range combinations {
		ok, err := g.luhnPkg.LuhnCheckCard(combination)
		if err != nil {
			return nil, err
		}
		if ok {
			validCards = append(validCards, combination)
		}
	}

	if pick {
		picked := combinations[rand.Intn(len(combinations))]
		return []string{picked}, nil
	}

	return validCards, nil
}

func (g *Generate) generateCombinations(card string, astNum int) []string {
	if astNum == 0 {
		return []string{card}
	}

	limit := g.intPow(10, astNum)
	cardTemplate := strings.Replace(card, "*", "%d", -1)
	combinations := make([]string, 0, limit)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			digits := fmt.Sprintf("%0*d", astNum, i)
			fullCard := fmt.Sprintf(cardTemplate, g.toInterfaceSlice(digits)...)
			mu.Lock()
			combinations = append(combinations, fullCard)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return combinations
}

func (g *Generate) intPow(base, exp int) int {
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}

func (g *Generate) toInterfaceSlice(s string) []interface{} {
	res := make([]interface{}, len(s))
	for i, r := range s {
		res[i] = r - '0'
	}
	return res
}

func (g *Generate) validateCardFormat(card string) error {
	if len(card) < 13 || len(card) > 19 {
		return ErrInvalidCardFormat
	}
	for i, c := range card {
		if c == '*' && i < len(card)-4 {
			return ErrInvalidCardFormat
		}
		if c != '*' && (c < '0' || c > '9') {
			return ErrInvalidCardFormat
		}
	}
	return nil
}
