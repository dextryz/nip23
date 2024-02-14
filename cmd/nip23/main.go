package main

import (
	"fmt"
	"os"

	"github.com/dextryz/nip23"
)

func main() {
	err := nip23.Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
