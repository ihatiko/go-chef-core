package go_chef_modules_sdk

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"runtime/debug"
)

type Module struct {
	MainCommand *cobra.Command
}

func NewModule() *Module {
	var versionCmd = &cobra.Command{
		Use:  "version",
		Long: "get current version module",
		Run: func(cmd *cobra.Command, args []string) {
			dbg, err := debug.ParseBuildInfo("")
			if err != nil {
				slog.Error("Error parsing build info", slog.String("error", err.Error()))
			}
			fmt.Println(dbg.Main.Version)
		},
	}
	mainCommand := &cobra.Command{}
	mainCommand.AddCommand(versionCmd)
	module := &Module{
		MainCommand: mainCommand,
	}
	return module
}
func (module *Module) AddCommands(commands ...*cobra.Command) {
	module.MainCommand.AddCommand(commands...)
}

func (module *Module) NewCommand(call string, fn func(cmd *cobra.Command, args []string)) *cobra.Command {
	return module.NewDetailCommand(call, "", "", fn)
}

func (module *Module) NewDetailCommand(use, short, long string, fn func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run:   fn,
	}
}

func (module *Module) Run() {
	if err := module.MainCommand.Execute(); err != nil {
		slog.Error("Error executing command", slog.String("error", err.Error()))
	}
}
