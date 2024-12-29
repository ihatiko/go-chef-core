package main

import (
	_ "embed"
	"errors"
	"fmt"
	gochefcodegenutils "github.com/ihatiko/go-chef-code-gen-utils"
	sdk "github.com/ihatiko/go-chef-modules-sdk"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "go-chef-core-modules.toml"
const configDir = ".go-chef"
const mainConfigName = "go-chef.toml"

type Module struct {
	Desc       string `toml:"desc"`
	Path       string `toml:"path"`
	Deprecated bool   `toml:"deprecated"`
}
type Config struct {
	Modules map[string]Module `toml:"modules"`
}
type MainConfig struct {
	Proxies []string `toml:"proxies"`
}

//go:embed config.toml
var defaultConfig []byte

func main() {
	cfg := getConfig()
	dir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting user cache dir", slog.Any("err", err))
		return
	}
	configPath := filepath.Join(dir, configDir, configFileName)
	module := sdk.NewModule()
	module.AddCommands(
		module.NewDetailCommand("show-config", "Show configuration", "", func(cmd *cobra.Command, args []string) {
			bytes, err := os.ReadFile(configPath)
			if err != nil {
				return
			}
			fmt.Println(string(bytes))
		}),
		module.NewDetailCommand("reset-config", "Reset configuration to base setup", "", func(cmd *cobra.Command, args []string) {
			err := os.WriteFile(configPath, defaultConfig, os.ModePerm)
			if err != nil {
				return
			}
			slog.Info("reset config", slog.String("path", configPath))
		}),
		module.NewDetailCommand("config-path", "Show config path", "", func(cmd *cobra.Command, args []string) {
			fmt.Println(configPath)
		}),
	)

	mainConfigPath := filepath.Join(dir, configDir, mainConfigName)
	if _, err := os.Stat(mainConfigPath); errors.Is(err, fs.ErrNotExist) {
		slog.Error("mainConfigPath does not exist", slog.String("mainConfigPath", mainConfigPath))
		return
	}

	mainConfigBytes, err := os.ReadFile(mainConfigPath)
	if err != nil {
		slog.Error("Error reading main config file", slog.String("mainConfigPath", mainConfigPath))
		return
	}
	mainConfig := new(MainConfig)
	err = toml.Unmarshal(mainConfigBytes, mainConfig)
	if err != nil {
		slog.Error("Error parsing main config file", slog.String("mainConfigPath", mainConfigPath), slog.Any("err", err))
		return
	}

	updater := gochefcodegenutils.NewUpdater(mainConfig.Proxies)
	composer := gochefcodegenutils.NewExecutor()
	for key, md := range cfg.Modules {
		cm := module.NewDetailCommand(key, md.Desc, "", func(cmd *cobra.Command, args []string) {
			updater.AutoUpdate(md.Path)
			splittedPath := strings.Split(md.Path, "/")
			coreCommand := splittedPath[len(splittedPath)-1]
			params := strings.Join(os.Args[2:], " ")
			proxyCommand := fmt.Sprintf("%s %s", coreCommand, params)
			result, err := composer.ExecDefaultCommand(proxyCommand)
			if err != nil {
				slog.Error("Error executing command: ", slog.Any("error", err), slog.String("command", params))
			}
			fmt.Println(result)
		})
		module.AddCommands(cm)
	}
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

func getConfig() *Config {
	dir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting user cache dir", slog.Any("err", err))
		return nil
	}
	configPath := filepath.Join(dir, configDir, configFileName)
	state, err := exists(configPath)
	if err != nil {
		slog.Error("Error checking if config file exists", slog.Any("err", err))
		return nil
	}
	if !state {
		dirPath := filepath.Join(dir, configDir)
		dirState, err := exists(dirPath)
		if err != nil {
			slog.Error("Error checking if config file exists", slog.Any("err", err))
			return nil
		}
		if !dirState {
			slog.Info("Config dir does not exist, creating", slog.Any("dir", dirPath))
			err := os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				slog.Error("Error creating config dir", slog.Any("err", err))
				return nil
			}
		}
		slog.Info("try create system config file")
		err = os.WriteFile(configPath, defaultConfig, os.ModePerm)
		if err != nil {
			slog.Error("Error creating config file", slog.Any("err", err))
			return nil
		}
		slog.Info("config file created", slog.Any("file", configPath))
	}
	cfgBytes, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("Error reading config file", slog.Any("err", err), slog.String("file", configPath))
		return nil
	}
	config := new(Config)
	err = toml.Unmarshal(cfgBytes, config)
	if err != nil {
		slog.Error("Error unmarshalling default config", slog.Any("error", err))
		return nil
	}
	return config
}
