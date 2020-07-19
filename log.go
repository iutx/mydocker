package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"os"
	"path"
)

func logContainer(containerName string) {
	logPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerName), container.ContainerLogFile)
	file, err := os.Open(logPath)
	if err != nil {
		log.Errorf("Open log file %v error: ", logPath, err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Read log %v error: ", logPath, err)
		return
	}
	fmt.Fprintf(os.Stdout, string(content))
}
