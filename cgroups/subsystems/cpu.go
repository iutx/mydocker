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
			// cpu.shares Default value 1024, if set 512, cpu use limit 50%.
			if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "cpu.shares"), []byte(res.CpuShare), 0664); err != nil {
				return fmt.Errorf("set cgroup cpu error: #{err}")
			}
		} else {
			return err
		}
	}
	return nil
}

func (c *CpuSubsystem) Apply(cGroupPath string, pid int) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		// Only one pid? if no, why don't use append instead of overflow.
		if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("cpu: pid %v add to cgroup %v errors: %v", pid, cGroupPath, err)
		}
		return nil
	} else {
		return err
	}
}

func (c *CpuSubsystem) Remove(cGroupPath string) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		return os.RemoveAll(subSysCGroupPath)
	} else {
		return err
	}
}
