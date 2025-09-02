package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseNetstatWindows(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock netstat script
	netstatScript := "#!/bin/sh\nprintf \"%s\" \"$MOCK_NETSTAT_OUTPUT\"\n"
	netstatPath := filepath.Join(tempDir, "netstat")
	if err := os.WriteFile(netstatPath, []byte(netstatScript), 0o755); err != nil {
		t.Fatalf("failed to write netstat script: %v", err)
	}

	// Create mock tasklist script
	tasklistScript := "#!/bin/sh\npid=${2##* }\ncat <<EOF\nImage Name                     PID Session Name        Session#    Mem Usage\n========================= ======== ================ =========== ============\n\nproc${pid}.exe                 ${pid} Console                    1     3,000 K\nEOF\n"
	tasklistPath := filepath.Join(tempDir, "tasklist")
	if err := os.WriteFile(tasklistPath, []byte(tasklistScript), 0o755); err != nil {
		t.Fatalf("failed to write tasklist script: %v", err)
	}

	// Prepend tempDir to PATH so our scripts are used
	origPath := os.Getenv("PATH")
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+origPath)

	tests := []struct {
		name   string
		output string
		want   []PortEntry
	}{
		{
			name:   "tcp",
			output: "TCP    127.0.0.1:80       0.0.0.0:0       LISTENING       1234\n",
			want: []PortEntry{{
				Protocol:      "TCP",
				LocalAddress:  "127.0.0.1:80",
				RemoteAddress: "0.0.0.0:0",
				State:         "LISTENING",
				PID:           "1234",
				ProcessName:   "proc1234.exe",
			}},
		},
		{
			name:   "udp",
			output: "UDP    127.0.0.1:53       *:*             5678\n",
			want: []PortEntry{{
				Protocol:      "UDP",
				LocalAddress:  "127.0.0.1:53",
				RemoteAddress: "*",
				State:         "N/A",
				PID:           "5678",
				ProcessName:   "proc5678.exe",
			}},
		},
		{
			name:   "both",
			output: "TCP    127.0.0.1:80       0.0.0.0:0       LISTENING       1234\nUDP    127.0.0.1:53       *:*             5678\n",
			want: []PortEntry{
				{
					Protocol:      "TCP",
					LocalAddress:  "127.0.0.1:80",
					RemoteAddress: "0.0.0.0:0",
					State:         "LISTENING",
					PID:           "1234",
					ProcessName:   "proc1234.exe",
				},
				{
					Protocol:      "UDP",
					LocalAddress:  "127.0.0.1:53",
					RemoteAddress: "*",
					State:         "N/A",
					PID:           "5678",
					ProcessName:   "proc5678.exe",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MOCK_NETSTAT_OUTPUT", tt.output)
			got, err := parseNetstatWindows()
			if err != nil {
				t.Fatalf("parseNetstatWindows() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseNetstatWindows() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
