package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
)

func Run(tty bool, command string, res *subsystems.ResourceConfig) {
	parent := container.NewParentProcess(tty, command)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

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
}
