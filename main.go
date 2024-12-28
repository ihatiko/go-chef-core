package main

import (
	"errors"
	"fmt"
	sdk "github.com/ihatiko/go-chef-modules-sdk"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
)

const fileName = "go-chef-core-modules.toml"
const configDir = ".go-chef"

type Module struct {
	Desc       string `toml:"desc"`
	Path       string `toml:"path"`
	Deprecated bool   `toml:"deprecated"`
}
type Modules map[string]Module

func main() {
	//dir, err := os.UserConfigDir()
	//if err != nil {
	//	slog.Error("Error getting user cache dir", slog.Any("err", err.Error()))
	//	return
	//}
	//dirPath := filepath.Join(dir, configDir)
	//if state, err := exists(dirPath); err != nil || state {
	//	err = os.Mkdir(dirPath, os.ModePerm)
	//	if err != nil {
	//		slog.Error("Error creating user cache dir", slog.Any("err", err.Error()))
	//		return
	//	}
	//}

	//fmt.Println(dir, err)
	module := sdk.NewModule()
	module.AddCommands(
		module.NewCommand("sandbox", func(cmd *cobra.Command, args []string) {
			fmt.Println("sandbox test")
		}),
		module.NewDetailCommand("show-config", "Show configuration", "", func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		}),
		module.NewDetailCommand("reset-config", "Reset configuration to base setup", "", func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		}),
		module.NewDetailCommand("config-path", "Show config path", "", func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO")
		}),
	)
	module.Run()
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
