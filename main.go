package main

import (
	"fmt"
	sdk "github.com/ihatiko/go-chef-proxy/go-chef-modules-sdk"
	"github.com/spf13/cobra"
)

func main() {
	module := sdk.NewModule()
	module.AddCommands(
		module.NewCommand("sandbox test", func(cmd *cobra.Command, args []string) {
			fmt.Println("sandbox test")
		}),
	)
	module.Run()
}
