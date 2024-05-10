package oneuptime_InfrastructureAgent_go

import (
	"github.com/gookit/slog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

func getCpuMetrics() *CPUMetrics {
	avg, err := load.Avg()
	if err != nil {
		slog.Error(err)
		return nil
	}

	numCpu, err := cpu.Counts(true)
	if err != nil {
		slog.Error(err)
		return nil
	}

	// Calculate CPU usage, which is the average load over the last minute divided by the number of CPUs
	cpuUsage := (avg.Load1 / float64(numCpu)) * 100
	return &CPUMetrics{
		PercentUsed: cpuUsage,
	}
}
