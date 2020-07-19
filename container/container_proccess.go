package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	RUNNING             string = "running"
	STOP                string = "stopped"
	Exit                string = "exited"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}


func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	// 创建匿名管道，获取 读取、写入 句柄
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Fatal("Pipe create error:", err)
	}

	// Must be use process clone from mydocker
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// 通过 EXTRA 携带管道读取句柄 创建子进程; 子进程为管道的读取端
	cmd.ExtraFiles = []*os.File{readPipe}
	imageURL := "/opt/busybox"
	// a index thing that is only needed for overlayfs do not totally
	// understand yet
	NewWorkspace(imageURL, volume)
	cmd.Dir = "/opt/merged"
	return cmd, writePipe
}

// Use default pipe. return write, read
func NewPipe() (*os.File, *os.File, error) {
	// 创建匿名管道
	if read, write, err := os.Pipe(); err != nil {
		return nil, nil, err
	} else {
		return read, write, nil
	}
}

func NewWorkspace(imageURL string, volume string) {
	mergedURL := "/opt/merged"
	indexURL := "/opt/index"
	writeLayerURL := "/opt/container_layer"

	if err := os.Mkdir(writeLayerURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeLayerURL, err)
	}
	if err := os.Mkdir(mergedURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mergedURL, err)
	}
	if err := os.Mkdir(indexURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", indexURL, err)
	}

	dirs := "lowerdir=" + imageURL + ",upperdir=" + writeLayerURL + ",workdir=" + indexURL
	log.Infof("overlayfs union parameters: %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}

	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			log.Infof("Volumes is: %v", volumeURLs)
			MountVolume(mergedURL, volumeURLs)
		} else {
			log.Errorf("Volume input param is not correct.")
		}
	}
}

func MountVolume(mergedURL string, volumeURLs []string) {
	hostPath := volumeURLs[0]
	if err := os.Mkdir(hostPath, 0777); err != nil {
		log.Info("Mkdir host path dir %s error: %v", hostPath, err)
	}
	containerPath := mergedURL + volumeURLs[1]
	if err := os.Mkdir(containerPath, 0777); err != nil {
		log.Info("Mkdir container path dir %s error: %v", containerPath, err)
	}

	cmd := exec.Command("mount", "--bind", hostPath, containerPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volumes error: ", err)
	}
}

func DeleteWorkSpace(volume string) {
	mergedURL := "/opt/merged"
	writeLayerURL := "/opt/container_layer"
	indexURL := "/opt/index"

	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			log.Infof("Volumes is: %v", volumeURLs)
			containerPath := mergedURL + volumeURLs[1]
			cmd := exec.Command("umount", containerPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Errorf("Umount volume failed. %v", err)
			}
		}
	}

	cmd := exec.Command("umount", mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}
	// remove merged, index and container write layer
	if err := os.RemoveAll(mergedURL); err != nil {
		log.Errorf("Remove dir %s error %v", mergedURL, err)
	}
	if err := os.RemoveAll(writeLayerURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeLayerURL, err)
	}
	if err := os.RemoveAll(indexURL); err != nil {
		log.Errorf("Remove dir %s error %v", indexURL, err)
	}
}
