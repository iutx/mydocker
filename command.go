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
		}, cli.StringFlag{
			Name:  "mem",
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
		var cmdArray []string

		for _, arg := range ctx.Args() {
			cmdArray = append(cmdArray, arg)
		}
		log.Infof("RunCommand command %s", cmdArray)
		// Judge have ti param.
		tty := ctx.Bool("ti")
		log.Infof("RunCommand tty bool %s", tty)
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}
		Run(tty, cmdArray, resConf)
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
}
