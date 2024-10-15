package main

import (
	"creditcard/internal/app"
	"fmt"
	"os"
)

func main() {
	a := app.NewApp()

	if err := a.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
