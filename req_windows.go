// +build windows

package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/alexbrainman/pc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

func addCounter(q *pc.Query, path string) *pc.Counter {
	c, err := q.AddCounter(path, 0)
	if err != nil {
		log.Println("addCounter: ", path, err)
		return nil
	}
	return c
}

func getFmtValueSafe(c *pc.Counter, format uint32) uint64 {
	var retval uint64
	_, rval, err := c.GetFmtValue(format)
	if err == nil {
		retval = rval.Value
	} else {
		log.Println("getFmtValueSafe", err)
	}
	return retval
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

func toFloat2(f float64) float64 {
	return math.Round(f*100) / 100
}

func getDiskInfo() map[string]float64 {
	diskMap := make(map[string]float64)
	parts, err := disk.Partitions(true)

	if err != nil {
		log.Println("get Partitions failed, err:%v\n", err)
	}

	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		diskMap[diskInfo.Path] = toFloat2(diskInfo.UsedPercent)
	}
	return diskMap
}

func ReqCounter(sqlInstance string) map[string]interface{} {
	const sleepMillisecond = 500

	q, err := pc.OpenQuery("", 0)
	if err != nil {
		log.Fatalln(err)
	}
	defer q.Close()

	mv := make(map[string]interface{})
	mc := make(map[string]*pc.Counter)

	c, _ := cpu.Percent(time.Millisecond*sleepMillisecond, false)
	mv["Processor"] = int(c[0])

	mc["Processor-pc"] = addCounter(q, `\Processor(_Total)\% Processor Time`)
	mc["ProcessorQueueLength"] = addCounter(q, `\System\Processor Queue Length`)
	mc["Memory"] = addCounter(q, `\Memory\% Committed Bytes In Use`)

	if len(sqlInstance) > 0 {
		mc["BatchRequests"] = addCounter(q, fmt.Sprintf(`\MSSQL$ %s:SQL Statistics\Batch Requests/sec`, sqlInstance))
		if mc["BatchRequests"] == nil {
			delete(mc, "BatchRequests")
		}

		mc["UserConnections"] = addCounter(q, fmt.Sprintf(`\MSSQL$ %s:General Statistics\User Connections`, sqlInstance))
		if mc["UserConnections"] == nil {
			delete(mc, "UserConnections")
		}
	}

	// resouce clean
	for _, m := range mc {
		defer m.Remove()
	}

	err = q.CollectData()
	if err != nil {
		log.Printf("CollectData failed: %v \n", err)
	}

	time.Sleep(time.Millisecond * sleepMillisecond)

	err = q.CollectData()
	if err != nil {
		log.Printf("CollectData failed: %v \n", err)
	}

	for k, m := range mc {
		mv[k] = getFmtValueSafe(m, pc.PDH_FMT_LARGE)
	}
	hostname, _ := os.Hostname()
	mv["HostsName"] = hostname
	mv["TimeStamp"] = time.Now().Format("2006-01-02 15:04:05")
	mv["IP"] = getIPAddr()
	mv["Disk"] = getDiskInfo()

	return mv
}
