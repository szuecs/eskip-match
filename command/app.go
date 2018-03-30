package command

import (
	"github.com/urfave/cli"
)

// NewApp creates the cli application
func NewApp() *cli.App {
	var configFile string
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configFile,
		},
	}
	app.Name = "eskip-match"
	app.Usage = "A command line tool that helps you test .eskip files routing matching logic"

	app.Commands = GetCommands(&Options{
		ConfigFile: configFile,
	})
	return app
}