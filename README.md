[![pipeline status](http://hqgit.inpark.kr/system-doc/go-perf/badges/master/pipeline.svg)](http://hqgit.inpark.kr/system-doc/go-perf/-/commits/master)

# go-perf
Performance mertic 수집 및 ElasticSearch 데이터 전송 

## Getting Started

### Prerequisites
- http://hqgit.inpark.kr/N15479/go-perf
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
2020/07/29 14:57:09 config : {"host":"http://localhost:9200","index":"perf","interval":5,"sqlinstance":""}
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
{"HostsName":"N15479","IP":"172.29.48.1","Memory":58,"Process":3,"ProcessorQueueLength":0,"TimeStamp":"2020-07-29 15:17:36"}

# Windows (w/SQLServer)
{"HostsName":"N15479","IP":"172.29.48.1","Memory":58,"Process":7,"ProcessorQueueLength":0,"TimeStamp":"2020-07-29 15:18:38","BatchRequests":0,"UserConnections":1}

# Posix 
{"cpu":0,"disk":1.17,"hostname":"N15479","ip":"172.28.26.123","loadavg":0.02,"mem":14.93,"swap":0,"timestamp":"2020-07-29 15:19:19"}
```
