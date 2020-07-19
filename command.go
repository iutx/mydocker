package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
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
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "Detach container.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Specified container name.",
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
		detach := ctx.Bool("d")

		if tty && detach {
			return fmt.Errorf("ti and d can't provided both")
		}

		log.Infof("RunCommand tty bool %s", tty)
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}
		volume := ctx.String("v")
		containerName := ctx.String("name")
		Run(tty, cmdArray, resConf, volume, containerName)
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

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit container into image.",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command.")
		}
		imageName := ctx.Args().Get(0)
		commitContainer(imageName)
		return nil
	},
}

var listCommand = cli.Command{
	Name:  "ps",
	Usage: "list all containers",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ps",
			Usage: "list all containers.",
		},
	},
	Action: func(ctx *cli.Context) error {
		ListContainers()
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "log",
	Usage: "show container log.",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command.")
		}
		containerName := ctx.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

var execCommand = cli.Command{
	Name :"exec",
	Usage: "Exec running container.",
	Action: func(ctx *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %s", os.Getgid())
			return nil
		}
		if len(ctx.Args()) < 2 {
			return fmt.Errorf("Missing container command.")
		}
		containerName := ctx.Args().Get(0)
		var commandArray []string
		for _, arg := range ctx.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
	},
}

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop container.",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command.")
		}
		containerName := ctx.Args().Get(0)
		stopContainer(containerName)
		return nil
	},
}
