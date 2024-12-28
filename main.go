package main

import (
	"errors"
	"fmt"
	sdk "github.com/ihatiko/go-chef-modules-sdk"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
)

const configDir = ".go-chef-core-modules"

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
		module.NewCommand("sandbox test", func(cmd *cobra.Command, args []string) {
			fmt.Println("sandbox test")
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
