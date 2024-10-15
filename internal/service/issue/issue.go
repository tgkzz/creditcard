package issue

import (
	"bufio"
	"creditcard/pkg/luhn"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Issuer interface {
	IssueCard(brandFile, issuerFile, brand, issuer string) (string, error)
}

type Issue struct {
	luhnPkg *luhn.Luhn
}

func NewIssuer() *Issue {
	return &Issue{
		luhnPkg: luhn.NewLuhn(),
	}
}

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
	ErrNotFound          = errors.New("not found")
)

func (is *Issue) IssueCard(brandFile, issuerFile, brand, issuer string) (string, error) {
	brandMap, err := is.loadBrands(brandFile)
	if err != nil {
		return "", err
	}

	if len(brandMap) == 0 {
		return "", ErrInvalidFileFormat
	}

	issuerMap, err := is.loadIssuer(issuerFile)
	if err != nil {
		return "", err
	}

	if len(issuerMap) == 0 {
		return "", ErrInvalidFileFormat
	}

	brandPr, err := is.getBrand(brand, brandMap)
	if err != nil {
		return "", err
	}

	issuerPr, err := is.getIssuer(issuer, issuerMap)
	if err != nil {
		return "", err
	}

	fullPrefix := brandPr + issuerPr

	left := 16 - len(fullPrefix)
	if left <= 1 {
		return "", ErrInvalidFileFormat
	}

	randDigit := make([]int, left)
	for i := 0; i < left-1; i++ {
		randDigit[i] = rand.Intn(10)
	}

	cardNumber := fullPrefix
	for _, digit := range randDigit {
		cardNumber += strconv.Itoa(digit)
	}

	checkDig := is.luhnPkg.CalculateLuhnCheckDigit(cardNumber)
	cardNumber += strconv.Itoa(checkDig)

	return cardNumber, nil
}

func (is *Issue) loadBrands(brands string) (map[string]string, error) {
	res := make(map[string]string)
	if err := is.loadDataToMap(brands, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (is *Issue) getBrand(brand string, mp map[string]string) (string, error) {
	for prefix, b := range mp {
		if b == brand {
			return prefix, nil
		}
	}
	return "", ErrNotFound
}

func (is *Issue) loadIssuer(issuers string) (map[string]string, error) {
	res := make(map[string]string)
	if err := is.loadDataToMap(issuers, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (is *Issue) getIssuer(issuer string, mp map[string]string) (string, error) {
	for prefix, k := range mp {
		if k == issuer {
			return prefix, nil
		}
	}
	return "", ErrNotFound
}

func (is *Issue) loadDataToMap(file string, mp map[string]string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return ErrInvalidFileFormat
		}
		mp[parts[1]] = parts[0]
	}
	return scanner.Err()

}
