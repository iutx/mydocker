package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	/*
	   1. CLONE_NEWUTS  UTS Namespace
	   2. CLONE_NEWIPC  IPC Namespace
	   3. CLONE_NEWPID  PID Namespace
	   4. CLONE_NEWNS   Mount Namespace
	   5. CLONE_NEWUSER User Namespace
	   6. CLONE_NEWNET  Network Namespace
	*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Run error:", err)
	}
}
