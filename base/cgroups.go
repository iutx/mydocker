package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {

	if os.Args[0] == "/proc/self/exe" {
		fmt.Println("Current PID:", syscall.Getpid())
		cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m --vm-keep -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("Run command err:", err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal("Run command err:", err)
		os.Exit(1)
	} else {
		fmt.Println("Outter pid:", cmd.Process.Pid)
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testMemoryLimit"), 0755)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testMemoryLimit", "tasks"),
			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testMemoryLimit", "memory.limit_in_bytes"),
			[]byte("100m"), 0644)
		cmd.Process.Wait()
	}
}
