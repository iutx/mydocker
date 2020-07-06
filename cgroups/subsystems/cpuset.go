package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSetSubsystem struct {
}

func (c *CpuSetSubsystem) Name() string {
	return "cpuset"
}

func (c *CpuSetSubsystem) Set(cGroupPath string, res *ResourceConfig) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset.cpus error: #{err}")
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}

func (c *CpuSetSubsystem) Apply(cGroupPath string, pid int) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		// Only one pid? if no, why don't use append instead of overflow.
		if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			// write error: No space left on device
			// https://blog.csdn.net/xftony/article/details/80536562
			return fmt.Errorf("pid %v add to cgroup %v errors: %v", pid, cGroupPath, err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}

func (c *CpuSetSubsystem) Remove(cGroupPath string) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		return os.RemoveAll(subSysCGroupPath)
	} else {
		return fmt.Errorf("get cgroup path error: %v\n", err)
	}
}
