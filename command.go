package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/container"
)


var runCommand = cli.Command{
	Name:  "run",
	Usage: "Run a container.",
	// Can use -- or -
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "start tty",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command.")
		}
		// command
		cmd := ctx.Args().Get(0)
		log.Infof("RunCommand command %s", cmd)
		// Judge have ti param.
		tty := ctx.Bool("ti")
		log.Infof("RunCommand tty bool %s", tty)
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "init container.",
	Action: func(ctx *cli.Context) error{
		cmd := ctx.Args().Get(0)
		log.Infof("InitCommand: command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

var testCommand = cli.Command{
	Name: "test",
	Usage: "test container.",
	Action: func(ctx *cli.Context) {
		cmd := ctx.Args().Get(0)
		log.Infof("testCommand: command %s", cmd)
	},
}