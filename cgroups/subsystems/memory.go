package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
}

func (m *MemorySubsystem) Name() string {
	return "memory"
}

func (m *MemorySubsystem) Set(cGroupPath string, res *ResourceConfig) error {
	if subSysCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory error: #{err}")
			}
		}
		return nil
	} else {
		return err
	}
}

func (m *MemorySubsystem) Apply(cGroupPath string, pid int) error {
	if subSysCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, false); err == nil {
		// Only one pid? if no, why don't use append instead of overflow.
		if err := ioutil.WriteFile(path.Join(subSysCGroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("memory: pid %v add to cgroup %v errors: %v", pid, cGroupPath, err)
		}
		return nil
	} else {
		return err
	}
}

func (m *MemorySubsystem) Remove(cGroupPath string) error {
	if subSysCGroupPath, err := GetCGroupPath(m.Name(), cGroupPath, false); err == nil {
		return os.RemoveAll(subSysCGroupPath)
	} else {
		return err
	}
}
