package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cligpt/shai/cmd"
)

func main() {
	ctx := context.Background()

	if err := cmd.Run(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
