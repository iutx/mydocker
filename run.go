package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"
)

func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	sendInitCommand(commandArray, writePipe)

	cGroupManager := cgroups.NewCGroupManager("mydocker-cgroup")
	// Destroy cGroup when exit container.
	defer cGroupManager.Destroy()

	if err := cGroupManager.Set(res); err != nil {
		log.Error("cGroup resource set error:", err)
	}
	if err := cGroupManager.Apply(parent.Process.Pid); err != nil {
		log.Error("cGroup resource applys error:", err)
	}

	if err := parent.Wait(); err != nil {
		log.Fatalf("Process wait error: ", err)
	}
	os.Exit(0)
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	log.Info("command all is:", command)
	if _, err := writePipe.WriteString(command); err != nil {
		log.Fatal("pipe write error:", err)
	}
	if err := writePipe.Close(); err != nil {
		log.Fatal("close pipe error:", err)
	}
}
