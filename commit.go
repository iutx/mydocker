package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func commitContainer(imageName string){
	mergedURL := "/opt/merged"
	imageTar := "/opt/" + imageName + ".tar"
	fmt.Printf("exported %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mergedURL, err)
	}
}
