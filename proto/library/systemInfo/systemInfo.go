package systemInfo

import (
	"context"
	"fmt"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/opq-osc/Yui/plugin/meta"
	"github.com/opq-osc/Yui/proto/library/systemInfo/export"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	PluginInfo meta.PluginMeta
}

func (s SystemInfo) CpuInfo(ctx context.Context, empty *emptypb.Empty) (*export.SystemInfoReply, error) {
	if s.PluginInfo.Permissions&meta.SystemInfoPermission == 0 {
		return nil, fmt.Errorf("缺少获取系统信息权限")
	}
	result, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuInfo := make([]export.CPUStat, len(result))
	for i := range result {
		cpuInfo[i] = export.CPUStat(result[i])
	}
	data, err := (&export.CpuInfo{CPU: cpuInfo}).MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &export.SystemInfoReply{Data: data}, nil
}
func (s SystemInfo) MemInfo(ctx context.Context, empty *emptypb.Empty) (*export.SystemInfoReply, error) {
	if s.PluginInfo.Permissions&meta.SystemInfoPermission == 0 {
		return nil, fmt.Errorf("缺少获取系统信息权限")
	}
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}
	data, err := (&export.MemInfo{VirtualMemory: (*struct {
		Total          uint64  `json:"total"`
		Available      uint64  `json:"available"`
		Used           uint64  `json:"used"`
		UsedPercent    float64 `json:"usedPercent"`
		Free           uint64  `json:"free"`
		Active         uint64  `json:"active"`
		Inactive       uint64  `json:"inactive"`
		Wired          uint64  `json:"wired"`
		Laundry        uint64  `json:"laundry"`
		Buffers        uint64  `json:"buffers"`
		Cached         uint64  `json:"cached"`
		WriteBack      uint64  `json:"writeBack"`
		Dirty          uint64  `json:"dirty"`
		WriteBackTmp   uint64  `json:"writeBackTmp"`
		Shared         uint64  `json:"shared"`
		Slab           uint64  `json:"slab"`
		Sreclaimable   uint64  `json:"sreclaimable"`
		Sunreclaim     uint64  `json:"sunreclaim"`
		PageTables     uint64  `json:"pageTables"`
		SwapCached     uint64  `json:"swapCached"`
		CommitLimit    uint64  `json:"commitLimit"`
		CommittedAS    uint64  `json:"committedAS"`
		HighTotal      uint64  `json:"highTotal"`
		HighFree       uint64  `json:"highFree"`
		LowTotal       uint64  `json:"lowTotal"`
		LowFree        uint64  `json:"lowFree"`
		SwapTotal      uint64  `json:"swapTotal"`
		SwapFree       uint64  `json:"swapFree"`
		Mapped         uint64  `json:"mapped"`
		VmallocTotal   uint64  `json:"vmallocTotal"`
		VmallocUsed    uint64  `json:"vmallocUsed"`
		VmallocChunk   uint64  `json:"vmallocChunk"`
		HugePagesTotal uint64  `json:"hugePagesTotal"`
		HugePagesFree  uint64  `json:"hugePagesFree"`
		HugePagesRsvd  uint64  `json:"hugePagesRsvd"`
		HugePagesSurp  uint64  `json:"hugePagesSurp"`
		HugePageSize   uint64  `json:"hugePageSize"`
	})(virtualMemory), SwapMemory: (*struct {
		Total       uint64  `json:"total"`
		Used        uint64  `json:"used"`
		Free        uint64  `json:"free"`
		UsedPercent float64 `json:"usedPercent"`
		Sin         uint64  `json:"sin"`
		Sout        uint64  `json:"sout"`
		PgIn        uint64  `json:"pgIn"`
		PgOut       uint64  `json:"pgOut"`
		PgFault     uint64  `json:"pgFault"`
		PgMajFault  uint64  `json:"pgMajFault"`
	})(swapMemory)}).MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &export.SystemInfoReply{Data: data}, nil
}
