package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mydocker/container"
	"os/exec"
)

func commitContainer(containerName string, imageName string) {
	mergedURL := fmt.Sprintf(container.DefaultMergedLocation, containerName)
	imageTar := fmt.Sprintf(container.DefaultImagesLocation, imageName)
	fmt.Printf("exported %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mergedURL, err)
	}
}
