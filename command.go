package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command.")
		}
		// command
		cmd := ctx.Args()
		log.Infof("RunCommand command %s", cmd)
		// Judge have ti param.
		tty := ctx.Bool("ti")
		log.Infof("RunCommand tty bool %s", tty)

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuShare:    ctx.String("cpushare"),
			CpuSet:      ctx.String("cpuset"),
		}

		Run(tty, cmd, resConf)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "init container.",
	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		log.Infof("InitCommand: command %s", cmd)
		err := container.RunContainerInitProcess()
		return err
	},
<<<<<<< HEAD
}

var testCommand = cli.Command{
	Name:  "test",
	Usage: "test container.",
	Action: func(ctx *cli.Context) {
		cmd := ctx.Args().Get(0)
		log.Infof("testCommand: command %s", cmd)
	},
}
=======
}
>>>>>>> 3.1
