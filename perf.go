package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type systemInfo struct {
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

func perf() systemInfo {
	chip := make(chan string)
	chmap := make(map[string](chan float64))
	chmap["cpu"] = make(chan float64)
	chmap["mem"] = make(chan float64)
	chmap["swap"] = make(chan float64)
	chmap["loadavg"] = make(chan float64)
	chmap["disk"] = make(chan float64)

	go func() {
		c, _ := cpu.Percent(time.Millisecond*300, false)
		chmap["cpu"] <- toFloat2(c[0])
	}()

	go func() {
		vm, _ := mem.VirtualMemory()
		chmap["mem"] <- toFloat2(vm.UsedPercent)

		swap, _ := mem.SwapMemory()
		chmap["swap"] <- toFloat2(swap.UsedPercent)
	}()

	go func() {
		disk, _ := disk.Usage("/")
		chmap["disk"] <- toFloat2(disk.UsedPercent)

		avg, _ := load.Avg()
		chmap["loadavg"] <- toFloat2(avg.Load5)
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

	var info = systemInfo{
		Timestamp: time.Now().Format("20060102150405"),
		Hostname:  hostname,
		Mem:       <-chmap["mem"],
		Cpu:       <-chmap["cpu"],
		Swap:      <-chmap["swap"],
		Disk:      <-chmap["disk"],
		Ip:        <-chip,
		LoadAvg:   <-chmap["loadavg"],
	}

	return info
}

func reqDo(idxName string, addrs []string) {
	info := perf()
	js, err := json.Marshal(info)
	if err != nil {
		fmt.Println("json.Marshal: %v", err)
		panic(err)
	}
	fmt.Println(string(js))

	cfg := elasticsearch.Config{Addresses: addrs}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("NewClient: %v", err)
		return
	}

	req := esapi.IndexRequest{
		Index:   idxName,
		Body:    strings.NewReader(string(js)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Println("req.Do: %v", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res)
	fmt.Println()
}
