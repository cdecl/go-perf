package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// perf = {
// 	#'@timestamp': ts,
// 	'timestamp': ts,
// 	'hostname': socket.gethostname(),
// 	'processor': psutil.cpu_percent(interval=1),
// 	'memory': psutil.virtual_memory().percent,
// 	'swap': psutil.swap_memory().percent,
// 	'disk': psutil.disk_usage('/').percent,
// 	'ip': socket.gethostbyname(socket.getfqdn())
// }
// 	perf["loadavg"] = psutil.getloadavg()[1]

type SystemInfo struct {
	Timestamp string  `json:"timestamp"`
	Hostname  string  `json:"hostname"`
	Cpu       float64 `json:"processor"`
	Mem       float64 `json:"memory"`
	Swap      float64 `json:"swap"`
	Disk      float64 `json:"disk"`
	Ip        string  `json:"ip"`
	LoadAvg   float64 `json:"loadavg"`
}

func toFloat2(f float64) float64 {
	return math.Round(f*100) / 100
}

func main() {
	chm := make(chan float64)
	chc := make(chan float64)
	chavg := make(chan float64)
	chswap := make(chan float64)
	chdisk := make(chan float64)
	chip := make(chan string)

	go func() {
		c, _ := cpu.Percent(time.Millisecond*300, false)
		chc <- toFloat2(c[0])
	}()

	go func() {
		vm, _ := mem.VirtualMemory()
		chm <- toFloat2(vm.UsedPercent)

		swap, _ := mem.SwapMemory()
		chswap <- toFloat2(swap.UsedPercent)
	}()

	go func() {
		avg, _ := load.Avg()
		chavg <- toFloat2(avg.Load5)
	}()

	go func() {
		disk, _ := disk.Usage("/")
		chdisk <- toFloat2(disk.UsedPercent)
	}()

	go func() {
		var ip string
		hostname, _ := os.Hostname()
		addrs, _ := net.LookupIP(hostname)
		for _, addr := range addrs {
			if addr.To4() != nil {
				ip = addr.To4().String()
				break
			}
		}
		chip <- ip
	}()

	hostname, _ := os.Hostname()

	var info = SystemInfo{
		Timestamp: time.Now().Format("20060102150405"),
		Hostname:  hostname,
		Mem:       <-chm,
		Cpu:       <-chc,
		Swap:      <-chswap,
		Disk:      <-chdisk,
		Ip:        <-chip,
		LoadAvg:   <-chavg,
	}

	js, _ := json.Marshal(info)
	fmt.Println(string(js))

}
