package main

import (
	"encoding/json"
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
	files, err := ioutil.ReadDir(infoPath)
	if err != nil {
		log.Errorf("Read dir %v error: %v", infoPath, err)
	}

	var containerInfos []*container.ContainerInfo
	for _, file := range files {
		configPath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, file.Name()), container.ConfigName)
		containerInfo, err := getContainerInfo(configPath)
		if err != nil {
			log.Errorf("Get container info error: ", err)
			return
		}
		containerInfos = append(containerInfos, containerInfo)
	}

	showInfo2Screen(containerInfos)
}

func getContainerInfo(configPath string) (*container.ContainerInfo, error) {
	data, _ := ioutil.ReadFile(configPath)
	containerInfo := container.ContainerInfo{}
	if err := json.Unmarshal(data, &containerInfo); err != nil {
		log.Errorf("Unmarshal file %v to json error: %v", configPath, err)
		return nil, err
	}

	return &containerInfo, nil
}

func showInfo2Screen(containerInfos []*container.ContainerInfo) {
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
