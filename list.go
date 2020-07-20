package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"os"
	"path"
	"text/tabwriter"
)

func ListContainers() {
	infoPath := fmt.Sprintf(container.DefaultInfoLocation, "")
	infoPath = infoPath[:len(infoPath)-1]
	exist, _ := container.PathExists(infoPath)
	if !exist {
		if err := os.MkdirAll(infoPath, 0622); err != nil {
			log.Errorf("Make dir %s error %v", infoPath, err)
			return
		}
	}
	files, err := ioutil.ReadDir(infoPath)
	if err != nil {
		log.Errorf("Read dir %v error: %v", infoPath, err)
	}

	var containerInfos []*container.ContainerInfo
	for _, file := range files {
		configPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, file.Name()), container.ConfigName)
		containerInfo, err := container.GetContainerInfo(configPath)
		if err != nil {
			log.Errorf("Get container info error: ", err)
			return
		}
		containerInfos = append(containerInfos, containerInfo)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containerInfos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreateTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}
