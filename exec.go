package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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

	// Get environment from container pid.
	containerEnvs := getEnvByPid(pid)
	// Add host environment and container environment to exec process.
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %v, error: %v", containerName, err)
	}
}

func getEnvByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error: %s", path, err)
	}
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}
