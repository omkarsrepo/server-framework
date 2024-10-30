package sf

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CommandsService interface {
	RegisterCommands()
}

type commandsService struct {
	cmd *cobra.Command
}

func NewCommandsService(command *cobra.Command) CommandsService {
	return &commandsService{
		cmd: command,
	}
}

func (props *commandsService) markRequiredFlags(flags ...string) {
	lo.ForEach(flags, func(flagName string, _ int) {
		_ = props.cmd.MarkPersistentFlagRequired(flagName)
	})
}

func (props *commandsService) register() {
	props.cmd.PersistentFlags().String("env", "sandbox", `env of server, can be "prod", "sandbox", "devl" or "localhost"`)
	props.cmd.PersistentFlags().Int("port", 8083, "port of server")
	props.cmd.PersistentFlags().Int("gracefulShutdownSecs", 1, "graceful shutdown secs for server")

	props.markRequiredFlags()
}

func (props *commandsService) bindToConfig() {
	viper.SetDefault("env", "sandbox")

	err := viper.BindPFlag("env", props.cmd.PersistentFlags().Lookup("env"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("port", props.cmd.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("gracefulShutdownSecs", props.cmd.PersistentFlags().Lookup("gracefulShutdownSecs"))
	if err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
}

func (props *commandsService) RegisterCommands() {
	props.register()
	props.bindToConfig()
}
