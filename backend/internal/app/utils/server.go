package utils

import (
	"backend/internal/app/global"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type Server struct {
	Os    Os      `json:"os"`
	Cpu   Cpu     `json:"cpu"`
	Ram   Ram     `json:"ram"`
	Disk  Disk    `json:"disk"`
	Cache []Cache `json:"cache"`
}

type Os struct {
	SysComputerName string `json:"sysComputerName"`
	SysOsName       string `json:"sysOsName"`
	SysOsArch       string `json:"sysOsArch"`
	NumCPU          int    `json:"numCpu"`
	Compiler        string `json:"compiler"`
	GoVersion       string `json:"goVersion"`
	NumGoroutine    int    `json:"numGoroutine"`
	RunTime         uint   `json:"runTime"`
}

type Cpu struct {
	CpuPercent float64 `json:"cpuPercent"`
	ModelName  string  `json:"modelName"`
	Cores      int     `json:"cores"`
}

type Ram struct {
	UsedMB      int     `json:"usedMb"`
	TotalMB     int     `json:"totalMb"`
	UsedPercent int     `json:"usedPercent"`
	SwapUsedMB  int     `json:"swapUsedMb"`
	SwapTotalMB int     `json:"swapTotalMb"`
	GoUsed      float64 `json:"goUsed"`
}

type Disk struct {
	UsedMB      int              `json:"usedMb"`
	UsedGB      int              `json:"usedGb"`
	TotalMB     int              `json:"totalMb"`
	TotalGB     int              `json:"totalGb"`
	UsedPercent int              `json:"usedPercent"`
	DiskList    []disk.UsageStat `json:"diskList"`
}

type Cache struct {
	bigcache.Stats
	Capacity int `json:"capacity"`
	Len      int `json:"len"`
}

// InitOS
// @description: OS信息
// @return: o Os, err error
func InitOS() (o Os) {
	sysInfo, err := host.Info()
	if err == nil {
		// 系统名称
		o.SysComputerName = sysInfo.Hostname
		// 操作系统
		o.SysOsName = sysInfo.OS
		// 系统架构
		o.SysOsArch = sysInfo.KernelArch
	}
	// 逻辑内核数
	o.NumCPU = runtime.NumCPU()
	o.Compiler = runtime.Compiler
	o.GoVersion = runtime.Version()
	o.NumGoroutine = runtime.NumGoroutine()
	// 程序运行时间
	// 计算开始和现在时间差
	o.RunTime = uint(time.Since(global.StartTime).Minutes())
	return o
}

// InitCPU
// @description: CPU信息
// @return: c Cpu, err error
func InitCPU() (c Cpu, err error) {
	// 物理内核数
	if cores, err := cpu.Counts(false); err != nil {
		return c, err
	} else {
		c.Cores = cores
	}
	// CPU型号
	if infos, err := cpu.Info(); err != nil {
		return c, err
	} else {
		c.ModelName = infos[0].ModelName
	}
	// CPU使用率
	if cpuPercent, err := cpu.Percent(time.Second, false); err != nil {
		return c, err
	} else {
		c.CpuPercent, _ = strconv.ParseFloat(strconv.FormatFloat(cpuPercent[0], 'f', 2, 64), 64)
	}
	return c, nil
}

// InitRAM
// @description: RAM信息
// @return: r Ram, err error
func InitRAM() (r Ram, err error) {
	// 物理内存
	if u, err := mem.VirtualMemory(); err != nil {
		return r, err
	} else {
		// 使用
		r.UsedMB = int(u.Used) / MB
		// 总量
		r.TotalMB = int(u.Total) / MB
		// 使用率
		r.UsedPercent = int(u.UsedPercent)

	}
	// 虚拟内存
	if u, err := mem.SwapMemory(); err != nil {
		return r, err
	} else {
		// 使用
		r.SwapUsedMB = int(u.Used) / MB
		// 总量
		r.SwapTotalMB = int(u.Total) / MB
	}

	// Golang占用内存
	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	r.GoUsed, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(gomem.Sys)/1024/1024), 64)

	return r, nil
}

// InitDisk
// @description: 硬盘信息
// @return: d Disk, err error
func InitDisk() (d Disk, err error) {
	diskInfo, err := disk.Partitions(true) //所有分区
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				d.DiskList = append(d.DiskList, *diskDetail)
			}
		}
	}

	return d, nil
}

func InitCache() (cache []Cache, err error) {
	cache = make([]Cache, 1)
	cache[1].Stats = global.TokenCache.Stats()
	cache[1].Len = global.TokenCache.Len()
	cache[1].Capacity = global.TokenCache.Capacity()
	return cache, err
}
