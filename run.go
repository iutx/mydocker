package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig, volume string, containerName string) {

	if containerName == "" {
		containerName = randStringBytes(6)
	}

	// 创建父进程以及管道写入句柄
	parent, writePipe := container.NewParentProcess(tty, volume, containerName)

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid, commandArray, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	cGroupManager := cgroups.NewCGroupManager("mydocker-cgroup")
	defer cGroupManager.Destroy()

	if err := cGroupManager.Set(res); err != nil {
		log.Errorf("cGroup %v set error.", cGroupManager.Path)
	}

	if err := cGroupManager.Apply(parent.Process.Pid); err != nil {
		log.Errorf("cGroup %v set error.", cGroupManager.Path)
	}

	sendInitCommand(commandArray, writePipe)

	if tty {
		if err := parent.Wait(); err != nil {
			log.Fatalf("Process wait error: ", err)
		}
		deleteContainerInfo(containerName)
	}
}

func deleteContainerInfo(containerName string) {
	containerInfoPath := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(containerInfoPath); err != nil {
		log.Errorf("Delete container info %v error: %v", containerInfoPath, err)
	}
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	log.Info("command all is:", command)
	// 通过管道写入句柄，写入执行命令到管道，等待子进程从管道中读取父进程传入的命令
	if _, err := writePipe.WriteString(command); err != nil {
		log.Fatal("pipe write error:", err)
	}
	if err := writePipe.Close(); err != nil {
		log.Fatal("close pipe error:", err)
	}
}

func recordContainerInfo(pid int, commandArray []string, containerName string) (string, error) {

	Id := randStringBytes(12)
	command := ""
	for _, cmd := range commandArray {
		command = command + " " + cmd
	}
	createTime := time.Now().Format("2006-01-02 15:04:05")

	containerInfo := container.ContainerInfo{
		Pid:        strconv.Itoa(pid),
		Id:         Id,
		Name:       containerName,
		Command:    command,
		CreateTime: createTime,
		Status:     container.RUNNING,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info err: %v", err)
		return "", err
	}

	configDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(configDir, 0622); err != nil {
		log.Errorf("Mkdir error: %v", err)
		return "", err
	}

	configPath := path.Join(configDir, container.ConfigName)
	file, err := os.Create(configPath)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %v error: %v", configPath, err)
		return "", err
	}

	if _, err := file.WriteString(string(jsonBytes)); err != nil {
		log.Errorf("File write string error: %v", err)
		return "", err
	}

	log.Infof("container info save success. path: %v", configPath)

	return containerName, nil
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
