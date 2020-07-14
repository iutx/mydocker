package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSetSubSystem struct {
}

func (c *CpuSetSubSystem) Name() string {
	return "cpuset"
}

func (c *CpuSetSubSystem) Set(cGroupPath string, res *ResourceConfig) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset error: #{err}")
			}
		}
		return nil
	} else {
		return err
	}
}

func (c *CpuSetSubSystem) Apply(cGroupPath string, pid int) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0664); err != nil {
			return fmt.Errorf("cpuset: pid %v add to cgroup %v errors: %v", pid, cGroupPath, err)
		}
		return nil
	} else {
		return err
	}
}

func (c *CpuSetSubSystem) Remove(cGroupPath string) error {
	if subSysCGroupPath, err := GetCGroupPath(c.Name(), cGroupPath, false); err == nil {
		return os.RemoveAll(subSysCGroupPath)
	} else {
		return err
	}
}
