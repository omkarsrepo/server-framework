package sfk

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

func (c *commandsService) markRequiredFlags(flags ...string) {
	lo.ForEach(flags, func(flagName string, _ int) {
		_ = c.cmd.MarkPersistentFlagRequired(flagName)
	})
}

func (c *commandsService) register() {
	c.cmd.PersistentFlags().String("env", "sandbox", `env of server, can be "prod", "sandbox", "devl" or "localhost"`)
	c.cmd.PersistentFlags().Int("port", 8083, "port of server")
	c.cmd.PersistentFlags().Int("gracefulShutdownSecs", 1, "graceful shutdown secs for server")

	c.markRequiredFlags()
}

func (c *commandsService) bindToConfig() {
	viper.SetDefault("env", "sandbox")

	err := viper.BindPFlag("env", c.cmd.PersistentFlags().Lookup("env"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("port", c.cmd.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag("gracefulShutdownSecs", c.cmd.PersistentFlags().Lookup("gracefulShutdownSecs"))
	if err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
}

func (c *commandsService) RegisterCommands() {
	c.register()
	c.bindToConfig()
}
