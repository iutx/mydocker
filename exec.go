package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/container"
	_ "mydocker/nsenter"
	"os"
	"os/exec"
	"strings"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func ExecContainer(containerName string, commandArray []string) {
	pid, err := container.GetPIDByContainerName(containerName)
	if err != nil {
		log.Errorf("Get PID erorr: %v", err)
		return
	}
	cmdStr := strings.Join(commandArray, " ")
	log.Infof("Current PID %v", pid)
	log.Infof("Command is %v", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %v, error: %v", containerName, err)
	}
}
