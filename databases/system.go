package databases

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var startTime = time.Now()

// SystemInfo represents the output for /api/server-info
type SystemInfo struct {
	Uptime           string `json:"uptime"`
	MemoryUsed       string `json:"memory_used"`
	MemoryTotal      string `json:"memory_total"`
	MemoryPercentage string `json:"memory_percentage"`
	CPUUsage         string `json:"cpu_usage"`
	OS               string `json:"os"`
	GoVersion        string `json:"go_version"`
	PostgresVersion  string `json:"postgres_version"`
}

// ServerInfoHandler mengembalikan status sistem real-time
func ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	usedMB, totalMB, memPct := getMemoryStats()

	info := SystemInfo{
		Uptime:           getUptime(),
		MemoryUsed:       fmt.Sprintf("%d MB", usedMB),
		MemoryTotal:      fmt.Sprintf("%d MB", totalMB),
		MemoryPercentage: memPct,
		CPUUsage:         getCPUUsage(),
		OS:               getOSName(),
		GoVersion:        runtime.Version(),
		PostgresVersion:  getPostgresVersion(),
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "sukses",
		"info":   info,
	})
}

func getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		// Fallback ke durasi uptime server-app Go
		uptimeSec := time.Since(startTime).Seconds()
		return formatSeconds(uptimeSec)
	}
	var sec float64
	fmt.Sscanf(string(data), "%f", &sec)
	return formatSeconds(sec)
}

func formatSeconds(sec float64) string {
	days := int(sec) / (24 * 3600)
	sec = sec - float64(days*24*3600)
	hours := int(sec) / 3600
	sec = sec - float64(hours*3600)
	minutes := int(sec) / 60
	seconds := int(sec) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d Hari", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d Jam", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d Menit", minutes))
	}
	parts = append(parts, fmt.Sprintf("%d Detik", seconds))

	if len(parts) > 2 {
		return strings.Join(parts[:2], " ")
	}
	return strings.Join(parts, " ")
}

func getMemoryStats() (usedMB, totalMB uint64, percentage string) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		// Fallback ke runtime Go mem stats
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		used := m.Alloc / 1024 / 1024
		sys := m.Sys / 1024 / 1024
		if sys == 0 {
			sys = 512
		}
		pct := fmt.Sprintf("%.1f%%", float64(used)/float64(sys)*100)
		return used, sys, pct
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var memTotal, memAvailable uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(line, "MemTotal: %d", &memTotal)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(line, "MemAvailable: %d", &memAvailable)
		}
	}

	if memTotal == 0 {
		return 0, 0, "0%"
	}

	usedKB := memTotal - memAvailable
	pct := fmt.Sprintf("%.1f%%", float64(usedKB)/float64(memTotal)*100)
	return usedKB / 1024, memTotal / 1024, pct
}

func getCPUUsage() string {
	t1, idle1, err := readCPUStats()
	if err != nil {
		return "0.0%"
	}
	time.Sleep(100 * time.Millisecond)
	t2, idle2, err := readCPUStats()
	if err != nil {
		return "0.0%"
	}
	totalDiff := float64(t2 - t1)
	idleDiff := float64(idle2 - idle1)
	if totalDiff == 0 {
		return "0.0%"
	}
	pct := (totalDiff - idleDiff) / totalDiff * 100.0
	return fmt.Sprintf("%.1f%%", pct)
}

func readCPUStats() (total, idle uint64, err error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[0] != "cpu" {
			return 0, 0, fmt.Errorf("invalid stat format")
		}

		var user, nice, system, idl, iowait, irq, softirq, steal uint64
		fmt.Sscanf(fields[1], "%d", &user)
		fmt.Sscanf(fields[2], "%d", &nice)
		fmt.Sscanf(fields[3], "%d", &system)
		fmt.Sscanf(fields[4], "%d", &idl)
		if len(fields) > 5 {
			fmt.Sscanf(fields[5], "%d", &iowait)
		}
		if len(fields) > 6 {
			fmt.Sscanf(fields[6], "%d", &irq)
		}
		if len(fields) > 7 {
			fmt.Sscanf(fields[7], "%d", &softirq)
		}
		if len(fields) > 8 {
			fmt.Sscanf(fields[8], "%d", &steal)
		}

		total = user + nice + system + idl + iowait + irq + softirq + steal
		idle = idl + iowait
		return total, idle, nil
	}
	return 0, 0, fmt.Errorf("empty stat file")
}

func getOSName() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				return strings.Trim(parts[1], `"`)
			}
		}
	}
	return "Linux"
}

func getPostgresVersion() string {
	dsn := getSuperadminDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "PostgreSQL"
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT version();").Scan(&version)
	if err != nil {
		return "PostgreSQL (Disconnect)"
	}

	fields := strings.Fields(version)
	if len(fields) >= 2 {
		return fields[0] + " " + fields[1]
	}
	return version
}
