package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

var (
	RUNNING                      string = "running"
	STOP                         string = "stopped"
	Exit                         string = "exited"
	DefaultInfoLocation          string = "/var/run/mydocker/%s/"
	ConfigName                   string = "config.json"
	ContainerLogFile             string = "container.log"
	BaseURL                      string = "/opt/mydocker"
	DefaultMergedLocation        string = BaseURL + "/merged/%s/"
	DefaultImagesLocation        string = BaseURL + "/images/%s.tar"
	DefaultReadOnlyLocationLayer string = BaseURL + "/base/%s"
	DefaulIndexLocation          string = BaseURL + "/index/%s"
	DefaultWritableLayerLocation string = BaseURL + "/container_layer/%s/"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	CreateTime  string   `json:"createTime"`
	Status      string   `json:"status"`
	Volume      string   `json:"volume"`
	PortMapping []string `json:"portmapping"`
}

func NewParentProcess(tty bool, volume string, containerName string, imageName string, envSlice []string) (*exec.Cmd, *os.File) {
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
	} else {
		containerDir := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(containerDir, 0622); err != nil {
			log.Errorf("NewParentProcess create file %v, error: %v", containerDir, err)
			return nil, nil
		}
		logPath := path.Join(containerDir, ContainerLogFile)
		if stdLogFile, err := os.Create(logPath); err == nil {
			cmd.Stdout = stdLogFile
		} else {
			log.Errorf("NewParentProcess create file %v, error: %v", logPath, err)
			return nil, nil
		}
	}

	// 通过 EXTRA 携带管道读取句柄 创建子进程; 子进程为管道的读取端
	cmd.ExtraFiles = []*os.File{readPipe}
	// Environment variable set.
	cmd.Env = append(os.Environ(), envSlice...)
	// a index thing that is only needed for overlayfs do not totally
	// understand yet
	NewWorkspace(containerName, volume, imageName)
	cmd.Dir = fmt.Sprintf(DefaultMergedLocation, containerName)
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

func createReadOnlyLayer(imageName string) error {
	unTarFolderUrl := fmt.Sprintf(DefaultReadOnlyLocationLayer, imageName)
	imageUrl := fmt.Sprintf(DefaultImagesLocation, imageName)
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", unTarFolderUrl, err)
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Errorf("Mkdir %s error %v", unTarFolderUrl, err)
			return err
		}

		log.Infof("Untar dir %s from %s", unTarFolderUrl, imageUrl)
		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			fmt.Println(imageUrl)
			log.Errorf("Untar dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

func createContainerLayer(mergedURL string, imageName string, indexURL string, writeLayerURL string) {
	if err := os.MkdirAll(writeLayerURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeLayerURL, err)
	}
	if err := os.MkdirAll(mergedURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mergedURL, err)
	}
	if err := os.MkdirAll(indexURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", indexURL, err)
	}
	baseURL := fmt.Sprintf(DefaultReadOnlyLocationLayer, imageName)

	dirs := "lowerdir=" + baseURL + ",upperdir=" + writeLayerURL + ",workdir=" + indexURL
	log.Infof("overlayfs union parameters: %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergedURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("%v", err)
	}

}

func NewWorkspace(containerName string, volume string, imageName string) {
	mergedURL := fmt.Sprintf(DefaultMergedLocation, containerName)
	indexURL := fmt.Sprintf(DefaulIndexLocation, containerName)
	writeLayerURL := fmt.Sprintf(DefaultWritableLayerLocation, containerName)

	if err := createReadOnlyLayer(imageName); err != nil {
		log.Errorf("Create readonly layer error: %v", err)
	}
	createContainerLayer(mergedURL, imageName, indexURL, writeLayerURL)

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

func DeleteWorkSpace(containerName string, volume string) {
	mergedURL := fmt.Sprintf(DefaultMergedLocation, containerName)
	indexURL := fmt.Sprintf(DefaulIndexLocation, containerName)
	writeLayerURL := fmt.Sprintf(DefaultWritableLayerLocation, containerName)

	if volume == "" {
		containerInfo, _ := GetContainerInfoByName(containerName)
		volume = containerInfo.Volume
	}

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

	containerPath := fmt.Sprintf(DefaultInfoLocation, containerName)

	if err := os.RemoveAll(containerPath); err != nil {
		log.Errorf("Remove dir %v error: %v", containerPath, err)
		return
	} else {
		log.Infof("Remove container %v success.", containerName)
	}
}
