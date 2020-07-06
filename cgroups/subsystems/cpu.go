package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubsystem struct {
	
}

func (c *CpuSubsystem) Name() string {
	return "cpu"
}

func (c *CpuSubsystem) Set(cGroupPath string, res *ResourceConfig) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, true); err == nil {
		if res.CpuShare != "" {
			if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu.shares error: #{err}")
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}

func (c *CpuSubsystem) Apply(cGroupPath string, pid int) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		// Only one pid? if no, why don't use append instead of overflow.
		if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("pid %v add to cgroup %v errors: %v", pid, cGroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}

func (c *CpuSubsystem) Remove(cGroupPath string) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		return os.RemoveAll(subSysCGroupPath)
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}
