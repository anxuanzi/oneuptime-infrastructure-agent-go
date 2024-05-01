package oneuptime_InfrastructureAgent_go

import (
	"github.com/shirou/gopsutil/v3/mem"
	"log/slog"
)

func getMemoryMetrics() *MemoryMetrics {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return &MemoryMetrics{
		Total:       memoryInfo.Total,
		Free:        memoryInfo.Free,
		Used:        memoryInfo.Used,
		PercentUsed: memoryInfo.UsedPercent,
		PercentFree: 100 - memoryInfo.UsedPercent,
	}
}
