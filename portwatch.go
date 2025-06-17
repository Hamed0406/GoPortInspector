package main

import (
	"fmt"
	"log"
	"os/exec"
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

func getProcessName(pid string) string {
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

func parseNetstat() ([]PortEntry, error) {
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
				name := getProcessName(pid)

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
				name := getProcessName(pid)

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

func main() {
	for {
		entries, err := parseNetstat()
		if err != nil {
			log.Printf("Error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Clear screen on Windows
		exec.Command("cmd", "/c", "cls").Run()

		fmt.Printf("%-6s %-22s %-22s %-12s %-6s %s\n", "Proto", "Local Address", "Remote Address", "State", "PID", "Process")
		fmt.Println(strings.Repeat("-", 90))
		for _, e := range entries {
			fmt.Printf("%-6s %-22s %-22s %-12s %-6s %s\n",
				e.Protocol, e.LocalAddress, e.RemoteAddress, e.State, e.PID, e.ProcessName)
		}

		time.Sleep(5 * time.Second)
	}
}
