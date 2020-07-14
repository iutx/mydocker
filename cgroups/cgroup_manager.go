package cgroups

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mydocker/cgroups/subsystems"
)

type CGroupManager struct {
	// CGroup path in hierarchy,
	Path     string
	Resource *subsystems.ResourceConfig
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Apply(pid int) error {
	for _, subsystem := range subsystems.SubSystemIns {
		if err := subsystem.Apply(c.Path, pid); err != nil {
			log.Fatalf("%v apply error: %v", subsystem.Name(), err)
		} else {
			log.Infof("%v apply success", subsystem.Name())
		}
	}
	return nil
}

func (c *CGroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subsystem := range subsystems.SubSystemIns {
		if err := subsystem.Set(c.Path, res); err != nil {
			return fmt.Errorf(" %v set error: %v", subsystem.Name(), err)
		} else {
			log.Infof(" %v set success.", subsystem.Name())
		}
	}
	return nil
}

func (c *CGroupManager) Destroy() error {
	for _, subsystem := range subsystems.SubSystemIns {
		if err := subsystem.Remove(c.Path); err != nil {
			log.Fatalf("%v remove error: %v", subsystem.Name(), err)
		} else {
			log.Infof("%v remove success", subsystem.Name())
		}
	}
	return nil
}
