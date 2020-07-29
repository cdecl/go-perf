// +build linux

package main

import (
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

func ReqCounter(sqlInstance string) map[string]interface{} {
	_ = sqlInstance
	mv := make(map[string]interface{})

	c, _ := cpu.Percent(time.Millisecond*300, false)
	mv["cpu"] = toFloat2(c[0])

	vm, _ := mem.VirtualMemory()
	mv["mem"] = toFloat2(vm.UsedPercent)

	swap, _ := mem.SwapMemory()
	mv["swap"] = toFloat2(swap.UsedPercent)

	disk, _ := disk.Usage("/")
	mv["disk"] = toFloat2(disk.UsedPercent)

	avg, _ := load.Avg()
	mv["loadavg"] = toFloat2(avg.Load5)
	mv["ip"] = getIPAddr()
	hostname, _ := os.Hostname()
	mv["hostname"] = hostname
	mv["timestamp"] = time.Now().Format("2006-01-02 15:04:05")

	return mv
}
