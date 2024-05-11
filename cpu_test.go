package oneuptime_InfrastructureAgent_go

import "testing"

func TestCpu(t *testing.T) {
	t.Log("Usage (%): ", getCpuMetrics().PercentUsed)
}
