package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {
	// 子进程读取父进程通过管道传递过来的命令
	commandArray := readUserCommand()
	log.Infof("RunContainerInitProcess command %s", commandArray)
	if commandArray == nil || len(commandArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	setUpMount()

	// SYS_EXECVE 系统调用不会在 path中寻找命令，通过 LookPath 寻找命令在系统中的绝对路径
	cmdAbsPath, err := exec.LookPath(commandArray[0])

	if err != nil {
		log.Errorf("Exec loop cmdAbsPath error %v", err)
		return err
	}

	// 通过寻找到的绝对路径，传入调用参数，环境信息；-》运行命令
	if err := syscall.Exec(cmdAbsPath, commandArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readUserCommand() []string {
	// 通过 index 3 去连接管道
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	// 读取父进程传入到管道中的命令
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

/**
Init 挂载点
*/
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)
	// https://github.com/xianlubird/mydocker/issues/41
	// 通过 mount /proc 挂载宿主机的 /proc ， ps 只能查看当前的容器内的进程就是通过这里实现。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	pivotRoot(pwd)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_PRIVATE|syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		log.Errorf("mount rootFS to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, "pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		log.Errorf("failed to create pivot root: %v", err)
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		log.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		log.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", "pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
