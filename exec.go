package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"os"
	"os/exec"
	"path"
	"strings"
	_ "mydocker/nsenter"
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func ExecContainer(containerName string, commandArray []string) {
	pid, err := getPIDByContainerName(containerName)
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

func getPIDByContainerName(containerName string) (string, error) {
	configPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerName), container.ConfigName)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Read container %v config error:  %v", containerName, err)
		return "", nil
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Unmarshal json error: %v", err)
		return "", err
	}
	return containerInfo.Pid, nil
}
