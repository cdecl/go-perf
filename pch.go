// +build windows

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/alexbrainman/pc"
)

func addCounter(q *pc.Query, path string) *pc.Counter {
	c, err := q.AddCounter(path, 0)
	if err != nil {
		fmt.Println(path, err)
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
		fmt.Println("getFmtValueSafe", err)
	}
	return retval
}

func ReqCounter(sqlInstance string) map[string]interface{} {
	q, err := pc.OpenQuery("", 0)
	if err != nil {
		fmt.Println(err)
	}
	defer q.Close()

	mv := make(map[string]interface{})
	mc := make(map[string]*pc.Counter)
	mc["Process"] = addCounter(q, `\Processor(_Total)\% Processor Time`)
	defer mc["Process"].Remove()

	mc["ProcessorQueueLength"] = addCounter(q, `\System\Processor Queue Length`)
	defer mc["ProcessorQueueLength"].Remove()

	mc["Memory"] = addCounter(q, `\Memory\% Committed Bytes In Use`)
	defer mc["Memory"].Remove()

	if len(sqlInstance) > 0 {
		mc["BatchRequests"] = addCounter(q, `\SQLServer:SQL Statistics\Batch Requests/sec`)
		defer mc["BatchRequests"].Remove()

		mc["UserConnections"] = addCounter(q, `\SQLServer:General Statistics\User Connections`)
		defer mc["UserConnections"].Remove()
	}

	err = q.CollectData()
	if err != nil {
		log.Printf("CollectData failed: %v \n", err)
	}

	const sleepMillisecond = 300
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

	getIPAddr := func() string {
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
	mv["IP"] = getIPAddr()

	return mv
}
