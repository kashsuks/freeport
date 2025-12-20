package api

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetBatteryPercentage() (int, error) {
	switch runtime.GOOS {
	case "darwin": // mac
		return getBatteryMacOS()
	case "linux":
		return getBatteryLinux()
	case "windows":
		return getBatteryWindows()
	default:
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func getBatteryMacOS() (int, error) {
	cmd := exec.Command("pmset", "-g", "batt")

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "%") {
			parts := strings.Split(line, "\t")
			if len(parts) >= 2 {
				percentStr := strings.TrimSpace(strings.Split(parts[1], ";")[0])
				percentStr = strings.TrimSuffix(percentStr, "%")
				return strconv.Atoi(percentStr)
			}
		}
	}

	return 0, fmt.Errorf("could not parse battery percentage")
}

func getBatteryLinux() (int, error) {
	cmd := exec.Command("cat", "/sys/class/power_supply/BAT0/capacity")

	output, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("cat", "/sys/class/power_supply/BAT1/capacity")
		output, err = cmd.Output()
		if err != nil {
			return 0, err
		}
	}

	percentStr := strings.TrimSpace(string(output))
	return strconv.Atoi(percentStr)
}

func getBatteryWindows() (int, error) {
	cmd := exec.Command("WMIC", "PATH", "Win32_Battery", "Get", "EstimatedChargeRemaining")
	
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		percentStr := strings.TrimSpace(lines[1])
		return strconv.Atoi(percentStr)
	}

	return 0, fmt.Errorf("could not parse battery percentage")
}