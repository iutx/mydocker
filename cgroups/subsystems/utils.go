package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func FindCGroupMountPoint(subsystem string) string {
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		field := strings.Split(line, " ")
		for _, subField := range strings.Split(field[len(field)-1], ",") {
			if subField == subsystem {
				return field[4]
			}
		}
	}
	return ""
}

func GetCGroupPath(subSystem string, cGroupPath string, autoCreate bool) (string, error) {
	cGroupRoot := FindCGroupMountPoint(subSystem) // subSystem param in struct is name field.
	cGroupAbsPath := path.Join(cGroupRoot, cGroupPath)
	// if cGroupAbsPath exist return path;
	// if cGroupAbsPath doesn't exist and allow create. create dir and return path.
	if _, err := os.Stat(cGroupAbsPath); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(cGroupAbsPath, 0755); err == nil {
			} else {
				return "", fmt.Errorf("mkdir error %v\n", cGroupAbsPath)
			}
		}
		return cGroupAbsPath, nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
