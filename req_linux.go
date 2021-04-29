// +build linux

package main

import (
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func toFloat2(f float64) float64 {
	return math.Round(f*100) / 100
}

func getIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func findStr(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func getDiskInfo() map[string]float64 {
	diskMap := make(map[string]float64)
	parts, err := disk.Partitions(true)
	ignore := []string{
		"/sys/fs/cgroup/devices", "binfmt_misc", "cgroup", "configfs", "debugfs", "devpts",
		"devtmpfs", "fusectl", "hugetlbfs", "mqueue", "nfsd", "nfs", "overlay", "proc", "pstore",
		"securityfs", "shm", "sunrpc", "sysfs", "systemd-1", "tmpfs", "autofs", "rpc_pipefs",
	}

	if err != nil {
		log.Println("get Partitions failed, err:%v\n", err)
	}

	for _, part := range parts {
		result := findStr(ignore, part.Fstype)
		if !result {
			diskInfo, _ := disk.Usage(part.Mountpoint)
			diskMap[diskInfo.Path] = toFloat2(diskInfo.UsedPercent)
		}
	}
	return diskMap
}

func ReqCounter(sqlInstance string) map[string]interface{} {
	_ = sqlInstance
	mv := make(map[string]interface{})

	c, _ := cpu.Percent(time.Millisecond*300, false)
	mv["processor"] = toFloat2(c[0])

	vm, _ := mem.VirtualMemory()
	mv["memory"] = toFloat2(vm.UsedPercent)

	swap, _ := mem.SwapMemory()
	mv["swap"] = toFloat2(swap.UsedPercent)

	// disk, _ := disk.Usage("/")
	mv["disk"] = getDiskInfo()

	avg, _ := load.Avg()
	mv["loadavg"] = toFloat2(avg.Load5)
	mv["ip"] = getIPAddr()
	hostname, _ := os.Hostname()
	mv["hostname"] = hostname
	mv["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return mv
}
