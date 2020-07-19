package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/container"
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

	container.DeleteWorkSpace(containerName, "")
}
