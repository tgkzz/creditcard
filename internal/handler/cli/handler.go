package cli

import (
	"creditcard/internal/service/generate"
	"creditcard/internal/service/info"
	"creditcard/internal/service/issue"
	"creditcard/internal/service/validate"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type IHandler interface {
	Gateway() error
}

type Handler struct {
	validator validate.Validator
	generator generate.Generator
	info      info.Info
	issuer    issue.Issuer
}

const (
	GENERATE    = "generate"
	VALIDATE    = "validate"
	INFORMATION = "information"
	ISSUE       = "issue"
)

var (
	ErrInvalidCommand     = errors.New("invalid command")
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrTooManyArguments   = errors.New("too many arguments")
)

func NewHandler() IHandler {
	return &Handler{
		validator: validate.NewValidator(),
		generator: generate.NewGenerator(),
		info:      info.NewInformation(),
		issuer:    issue.NewIssuer(),
	}
}

func (h *Handler) Gateway() error {
	args := os.Args[1:]

	c, err := h.getCommand(args)
	if err != nil {
		return err
	}

	switch c {
	case GENERATE:
		return h.generateHandler(args[1:])
	case VALIDATE:
		return h.validateHandler(args[1:])
	case INFORMATION:
		return h.infoHandler(args[1:])
	case ISSUE:
		return h.issueHandler(args[1:])
	default:
		return ErrInvalidCommand
	}
}

// TODO: generate is bagged, fix it
func (h *Handler) generateHandler(args []string) error {
	if len(args) == 0 {
		return ErrNotEnoughArguments
	}

	isPick := flag.Bool("pick", false, "randomly pick generated card number")

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	card := ""
	for _, arg := range args {
		if !strings.Contains(arg, "-") {
			card = arg
			break
		}
	}

	res, err := h.generator.Generate(card, *isPick)
	if err != nil {
		fmt.Println("echo $?")
		return err
	}

	for _, r := range res {
		fmt.Println(r)
	}
	return nil
}

func (h *Handler) validateHandler(args []string) error {
	for _, arg := range args {
		if err := h.validator.ValidateCard(arg); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (h *Handler) infoHandler(args []string) error {

	brands := flag.String("brands", "", "path to brands txt")
	issuers := flag.String("issuers", "", "path to issuer txt")
	useStdin := flag.Bool("stdin", false, "Read input from stdin")

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	cardNum := ""
	if *useStdin {
		if _, err := fmt.Scanln(&cardNum); err != nil {
			return err
		}
	} else {
		for _, arg := range args {
			if h.isOnlyNums(arg) {
				cardNum = arg
				break
			}
		}
	}

	if err := h.info.GetCardInfo(*brands, *issuers, cardNum); err != nil {
		return err
	}

	return nil
}

func (h *Handler) issueHandler(args []string) error {
	bFile := flag.String("brands", "", "Path to brands file")
	iFile := flag.String("issuers", "", "Path to issuer file")
	brand := flag.String("brand", "", "Card brand")
	issuer := flag.String("issuer", "", "Card issuer")

	if err := flag.CommandLine.Parse(args); err != nil {
		return err
	}

	if *brand == "" ||
		*issuer == "" ||
		*bFile == "" ||
		*iFile == "" {
		return errors.New("missing required parameters")
	}

	res, err := h.issuer.IssueCard(*bFile, *iFile, *brand, *issuer)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func (h *Handler) getCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", ErrNotEnoughArguments
	}

	cmd := args[0]

	return cmd, nil
}

func (h *Handler) isOnlyNums(str string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	if re.MatchString(str) {
		return true
	}

	return false
}
