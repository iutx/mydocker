package container

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

func RecordContainerInfo(pid int, commandArray []string, containerName string, volume string,
	portMapping []string) (*ContainerInfo, error) {

	Id := RandStringBytes(12)
	command := ""
	for _, cmd := range commandArray {
		command = command + " " + cmd
	}
	createTime := time.Now().Format("2006-01-02 15:04:05")

	containerInfo := &ContainerInfo{
		Pid:        strconv.Itoa(pid),
		Id:         Id,
		Name:       containerName,
		Command:    command,
		CreateTime: createTime,
		Status:     RUNNING,
		Volume:     volume,
		PortMapping: portMapping,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info err: %v", err)
		return nil, err
	}

	configDir := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(configDir, 0622); err != nil {
		log.Errorf("Mkdir error: %v", err)
		return nil, err
	}

	configPath := path.Join(configDir, ConfigName)
	file, err := os.Create(configPath)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %v error: %v", configPath, err)
		return nil, err
	}

	if _, err := file.WriteString(string(jsonBytes)); err != nil {
		log.Errorf("File write string error: %v", err)
		return nil, err
	}

	log.Infof("container info save success. path: %v", configPath)

	return containerInfo, nil
}

func GetPIDByContainerName(containerName string) (string, error) {
	configPath := path.Join(fmt.Sprintf(DefaultInfoLocation, containerName), ConfigName)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Read container %v config error:  %v", containerName, err)
		return "", nil
	}
	var containerInfo ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Unmarshal json error: %v", err)
		return "", err
	}
	return containerInfo.Pid, nil
}

func GetContainerInfoByName(containerName string) (*ContainerInfo, error) {
	configPath := path.Join(fmt.Sprintf(DefaultInfoLocation, containerName), ConfigName)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Read container %v config error:  %v", containerName, err)
		return nil, nil
	}
	var containerInfo ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Unmarshal json error: %v", err)
		return nil, err
	}
	return &containerInfo, nil
}

func RandStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GetContainerInfo(configPath string) (*ContainerInfo, error) {
	data, _ := ioutil.ReadFile(configPath)
	containerInfo := ContainerInfo{}
	if err := json.Unmarshal(data, &containerInfo); err != nil {
		log.Errorf("Unmarshal file %v to json error: %v", configPath, err)
		return nil, err
	}

	return &containerInfo, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
