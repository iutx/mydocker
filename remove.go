package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mydocker/container"
	"os"
)

func removeContainer(containerName string) {

	containerInfo, err := getContainerINfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %v info error: %v", containerName, err)
		return
	}
	if containerInfo.Status == container.RUNNING {
		log.Errorf("Please stop running container %v before remove it.", containerName)
		return
	}

	containerPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)

	if err := os.RemoveAll(containerPath); err != nil {
		log.Errorf("Remove dir %v error: %v", containerPath, err)
		return
	}else{
		log.Infof("Remove container %v success.", containerName)
	}
}
