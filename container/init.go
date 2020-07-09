package container

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func RunContainerInitProcess() error {
	// 子进程读取父进程通过管道传递过来的命令
	commandArray := readUserCommand()
	log.Infof("RunContainerInitProcess command %s", commandArray)

	// https://github.com/xianlubird/mydocker/issues/41
	// 通过 mount /proc 挂载宿主机的 /proc ， ps 只能查看当前的容器内的进程就是通过这里实现。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

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
