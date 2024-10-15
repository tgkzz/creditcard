package info

import (
	"bufio"
	"creditcard/pkg/luhn"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Info interface {
	GetCardInfo(brands, issuer, cardNum string) error
}

type Information struct {
	luhnPkg *luhn.Luhn
}

func NewInformation() Info {
	return &Information{
		luhnPkg: luhn.NewLuhn(),
	}
}

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
	ErrInvalidCardFormat = errors.New("invalid card format")
)

type Res struct {
	Card   string
	Valid  bool
	Brand  string
	Issuer string
}

func (i *Information) GetCardInfo(brands, issuer, cardNum string) error {
	brandsMap, err := i.loadBrands(brands)
	if err != nil {
		return err
	}

	issuersMap, err := i.loadIssuer(issuer)
	if err != nil {
		return err
	}

	res, err := i.getResInfo(cardNum, brandsMap, issuersMap)
	if err != nil {
		return err
	}

	i.printCardInfo(res)
	return nil
}

func (i *Information) printCardInfo(res *Res) {
	fmt.Println(res.Card)
	if res.Valid {
		fmt.Println("Correct: yes")
	} else {
		fmt.Println("Correct: no")
	}
	fmt.Printf("Card Brand: %s\n", res.Brand)
	fmt.Printf("Card Issuer: %s\n", res.Issuer)
}

func (i *Information) getResInfo(card string, brandsMap, issuerMap map[string]string) (*Res, error) {
	valid, err := i.luhnPkg.LuhnCheckCard(card)
	if err != nil {
		return nil, ErrInvalidCardFormat
	}

	return &Res{
		Card:   card,
		Valid:  valid,
		Brand:  i.getBrand(card, brandsMap),
		Issuer: i.getIssuer(card, issuerMap),
	}, nil
}

func (i *Information) loadBrands(brands string) (map[string]string, error) {
	res := make(map[string]string)
	if err := i.loadDataToMap(brands, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (i *Information) getBrand(card string, mp map[string]string) string {
	for prefix, brand := range mp {
		if strings.HasPrefix(card, prefix) {
			return brand
		}
	}
	return "-"
}

func (i *Information) loadIssuer(issuers string) (map[string]string, error) {
	res := make(map[string]string)
	if err := i.loadDataToMap(issuers, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (i *Information) getIssuer(card string, mp map[string]string) string {
	for prefix, issuer := range mp {
		if strings.HasPrefix(card, prefix) {
			return issuer
		}
	}
	return "-"
}

func (i *Information) loadDataToMap(file string, mp map[string]string) error {
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
