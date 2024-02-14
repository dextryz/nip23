package main

import (
	"fmt"
	"os"

	"github.com/dextryz/zet"
)

func main() {
	err := zet.Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
