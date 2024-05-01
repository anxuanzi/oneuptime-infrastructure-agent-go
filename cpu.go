package oneuptime_InfrastructureAgent_go

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"log/slog"
)

func getCpuMetrics() *CPUMetrics {
	avg, err := load.Avg()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	numCpu, err := cpu.Counts(true)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	cpuUsage := (avg.Load1 / float64(numCpu)) * 100
	return &CPUMetrics{
		PercentUsed: cpuUsage,
	}
}
