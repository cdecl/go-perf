[![pipeline status](http://hqgit.inpark.kr/system-doc/go-perf/badges/master/pipeline.svg)](http://hqgit.inpark.kr/system-doc/go-perf/-/commits/master)

# go-perf
Performance mertic 수집 및 ElasticSearch 데이터 전송 

## Getting Started

### Prerequisites
- http://hqgit.inpark.kr/util/go-perf
- golang : `go version go1.14.5 windows/amd64`

#### 수집기 
- psutil : https://github.com/shirou/gopsutil
  - golang 버전 psutil library 
  - posix metric 수집 (windows 호환가능)
- performance counter : https://github.com/alexbrainman/pc
  - windows performance counter helper api wrapper 

#### 서비스관리 
- https://github.com/kardianos/service
  - Service will install / un-install, start / stop, and run a program as a service (daemon)
  - Currently supports Windows XP+, Linux/(systemd | Upstart | SysV), and OSX/Launchd.

### Build

```sh
go mod tidy 
go build 
```

## Running the tests

### Config
- host: ElasticSearch Host 
  - If the value is empty, only console output
- index: ElasticSearch Index 
- interval: 수집간격(sec)
- sqlinstance: SQLServer Instance Name
  - Windows only, If the value is empty skip

```json
{
	"host": "http://localhost:9200",
	"index": "perf",
	"interval": 5,
	"sqlinstance": ""
}
```

### Run Test

```sh
$ ./go-perf 
2021/04/22 11:19:39 config : {"host":"http://localhost:9200","index":"perf","interval":20,"sqlinstance":""}
2021/04/22 11:19:39 START
2021/04/22 11:19:39 {"disk":{"/":11.76,"/boot":13.91,"/boot/efi":5.59,"/home":2.95,"/sys/firmware/efi/efivars":0},"hostname":"centos-cdecl","ip":"192.168.137.100","loadavg":0.01,"memory":19.13,"processor":0,"swap":0,"timestamp":"2021-04-22 11:19:39"}
2020/07/29 14:57:09 START
2020/07/29 14:57:10 {"HostsName":"N15479","IP":"172.29.48.1","Memory":57,"Process":9,"ProcessorQueueLength":3,"TimeStamp":"2020-07-29 14:57:10"}
2020/07/29 14:57:10 [201 Created] {"_index":"perf-20200729","_type":"_doc","_id":"nnklmXMBCjdbPcrUavha","_version":1,"result":"created","forced_refresh":true,"_shards":{"total":2,"successful":1,"failed":0},"_seq_no":242,"_primary_term":1}
...
```

### Install, Run
- 설치/삭제
  - windows : service 등록 
  - posix (init.d) : /etc/init.d/GoPerf
  - posix (systemd) : /etc/systemd/system/GoPerf.service

```sh
$ ./go-perf install 
$ ./go-perf uninstall 
```

- 시작/중지
  - windows : sc start GoPerf / sc stop GoPerf 
  - posix (init.d) : service GoPerf start / service GoPerf stop 
  - posix (systemd) : systemctl start GoPerf / systemctl stop GoPerf

```sh
$ ./go-perf start 
$ ./go-perf uninstall 
```

```sh
# Windows 
{"Disk":{"C:":91,"D:":13.06},"HostsName":"N15479-W02","IP":"192.168.137.1","Memory":86,"Processor":6,"Processor-pc":7,"ProcessorQueueLength":0,"TimeStamp":"2021-04-22 11:19:15"}

# Windows (w/SQLServer)
 {"BatchRequests":0,"Disk":{"C:":91,"D:":13.06},"HostsName":"N15479-W02","IP":"192.168.137.1","Memory":86,"Processor":4,"Processor-pc":4,"ProcessorQueueLength":0,"TimeStamp":"2021-04-22 11:18:40","UserConnections":0}

# Linux
{"disk":{"/":11.76,"/boot":13.91,"/boot/efi":5.59,"/home":2.95,"/sys/firmware/efi/efivars":0},"hostname":"centos-cdecl","ip":"192.168.137.100","loadavg":0.01,"memory":19.13,"processor":0,"swap":0,"timestamp":"2021-04-22 11:19:39"}
```
