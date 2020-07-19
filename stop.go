package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"path"
	"strconv"
	"syscall"
)

func stopContainer(containerName string) {
	pid, err := getPIDByContainerName(containerName)
	if err != nil {
		log.Errorf("Container %v get pid erorr: %v", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("PID %v convert error: %v", pid, err)
		return
	}

	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("PID %v kill error: %v", pidInt, err)
	}

	containerInfo, err := getContainerINfoByName(containerName)
	if err != nil {
		log.Errorf("Get container info error: %v", err)
	}

	containerInfo.Status = container.Exit
	containerInfo.Pid = ""

	newContent, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %v error: %v", containerName, err)
	}
	configPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerName), container.ConfigName)

	if err := ioutil.WriteFile(configPath, newContent, 0662); err != nil {
		log.Errorf("Container %v info write to config error: %v", containerName, err)
	}
}

func getContainerINfoByName(containerName string) (*container.ContainerInfo, error) {
	configPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerName), container.ConfigName)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Read container %v config error:  %v", containerName, err)
		return nil, nil
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Unmarshal json error: %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
