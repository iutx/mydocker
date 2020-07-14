package container

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
	cmd.Dir = "/root/busybox"
	// 通过 EXTRA 携带管道读取句柄 创建子进程; 子进程为管道的读取端
	cmd.ExtraFiles = []*os.File{readPipe}
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
