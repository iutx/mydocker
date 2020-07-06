package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string
	Set(cGroupPath string, res *ResourceConfig) error
	Apply(cGroupPath string, pid int) error
	Remove(cGroupPath string) error
}

var (
	SubSystemIns = []Subsystem{
		&CpuSubsystem{},
		&CpuSetSubsystem{},
		&MemorySubsystem{},
	}
)
