package main

import (
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"
)

func Run(tty bool, commandArray []string, res *subsystems.ResourceConfig) {
	// 创建父进程以及管道写入句柄
	parent, writePipe := container.NewParentProcess(tty)

	if err := parent.Start(); err != nil {
		log.Error(err)
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

	if err := parent.Wait(); err != nil {
		log.Fatalf("Process wait error: ", err)
	}
	mntURL := "/opt/mnt/"
	rootURL := "/opt/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
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
