package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PortEntry struct {
	Protocol      string
	LocalAddress  string
	RemoteAddress string
	State         string
	PID           string
	ProcessName   string
}

// getPortEntries returns port information using an OS specific implementation.
func getPortEntries() ([]PortEntry, error) {
	switch runtime.GOOS {
	case "windows":
		return parseNetstatWindows()
	case "linux", "darwin":
		return parseLsof()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func getProcessNameWindows(pid string) string {
	out, err := exec.Command("tasklist", "/FI", "PID eq "+pid).Output()
	if err != nil {
		return "Unknown"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) >= 4 {
		fields := strings.Fields(lines[3])
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return "Unknown"
}

func parseNetstatWindows() ([]PortEntry, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run netstat: %v", err)
	}

	var results []PortEntry
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TCP") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				proto := fields[0]
				local := fields[1]
				remote := fields[2]
				state := fields[3]
				pid := fields[4]
				name := getProcessNameWindows(pid)

				results = append(results, PortEntry{
					Protocol:      proto,
					LocalAddress:  local,
					RemoteAddress: remote,
					State:         state,
					PID:           pid,
					ProcessName:   name,
				})
			}
		} else if strings.HasPrefix(line, "UDP") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				proto := fields[0]
				local := fields[1]
				pid := fields[3]
				name := getProcessNameWindows(pid)

				results = append(results, PortEntry{
					Protocol:      proto,
					LocalAddress:  local,
					RemoteAddress: "*",
					State:         "N/A",
					PID:           pid,
					ProcessName:   name,
				})
			}
		}
	}
	return results, nil
}

// parseLsof parses lsof output for Linux and macOS.
func parseLsof() ([]PortEntry, error) {
	cmd := exec.Command("lsof", "-nP", "-i")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsof: %v", err)
	}

	var results []PortEntry
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		proto := fields[7]
		pid := fields[1]
		proc := fields[0]
		nameField := strings.Join(fields[8:], " ")

		local := ""
		remote := "*"
		state := "N/A"

		parts := strings.Split(nameField, "->")
		local = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			remoteParts := strings.Split(parts[1], " ")
			if len(remoteParts) > 0 {
				remote = remoteParts[0]
			}
		}
		if idx := strings.LastIndex(nameField, "("); idx != -1 {
			state = strings.TrimSuffix(nameField[idx+1:], ")")
		}

		results = append(results, PortEntry{
			Protocol:      proto,
			LocalAddress:  local,
			RemoteAddress: remote,
			State:         state,
			PID:           pid,
			ProcessName:   proc,
		})
	}
	return results, nil
}

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "cls").Run()
	default:
		exec.Command("clear").Run()
	}
}

func main() {
	for {
		entries, err := getPortEntries()
		if err != nil {
			log.Printf("Error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		clearScreen()

		fmt.Printf("%-6s %-22s %-22s %-12s %-6s %s\n", "Proto", "Local Address", "Remote Address", "State", "PID", "Process")
		fmt.Println(strings.Repeat("-", 90))
		for _, e := range entries {
			fmt.Printf("%-6s %-22s %-22s %-12s %-6s %s\n",
				e.Protocol, e.LocalAddress, e.RemoteAddress, e.State, e.PID, e.ProcessName)
		}

		time.Sleep(5 * time.Second)
	}
}
