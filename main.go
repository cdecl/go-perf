package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kardianos/service"
)

type Args struct {
	Host        interface{} `json:"host"`
	Index       string      `json:"index"`
	Interval    int64       `json:"interval"`
	SqlInstance string      `json:"sqlinstance"`
}

type program struct {
	args Args
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	p.reqDo()

	ticker := time.NewTicker(time.Second * time.Duration(p.args.Interval))
	for {
		select {
		case tm := <-ticker.C:
			_ = tm
			p.reqDo()
		case <-p.exit:
			ticker.Stop()
			return
		}
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	return nil
}

func (p *program) reqDo() {
	dicPerf := ReqCounter(p.args.SqlInstance)
	sb, err := json.Marshal(dicPerf)
	if err != nil {
		log.Println("ReqCounter: %v", err)
		return
	}

	hosts := getHost(p.args.Host)

	log.Println(string(sb))
	if len(hosts) == 0 {
		return
	}

	cfg := elasticsearch.Config{Addresses: hosts}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("NewClient: %v", err)
		return
	}

	req := esapi.IndexRequest{
		Index:   p.args.Index + "-" + time.Now().Format("20060102"),
		Body:    strings.NewReader(string(sb)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Println("reqDo: %v", err)
		return
	}
	defer res.Body.Close()
	log.Println(res)
}

func getConfigPath() string {
	fullexecpath, err := os.Executable()
	if err != nil {
		return ""
	}

	dir, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]

	return filepath.Join(dir, name+".json")
}

func getHost(h interface{}) []string {
	hosts := []string{}
	switch v := h.(type) {
	case []interface{}:
		for _, val := range h.([]interface{}) {
			hosts = append(hosts, val.(string))
		}
	case string:
		if len(h.(string)) > 0 {
			hosts = append(hosts, h.(string))
		}
	default:
		_ = v
	}
	return hosts
}

func getArgs() (Args, error) {
	args := Args{}

	f, err := os.Open(getConfigPath())
	if err != nil {
		return args, err
	}
	defer f.Close()

	r := json.NewDecoder(f)
	err = r.Decode(&args)
	if err != nil {
		return args, err
	}
	return args, nil
}

func main() {
	args, err := getArgs()
	if err != nil {
		log.Fatal("config load error")
		return
	}

	js, _ := json.Marshal(args)
	log.Printf("config : %v", string(js))
	log.Printf("START")

	svcConfig := &service.Config{
		Name:        "GoPerf",
		Description: "INTERPARK GoPerf",
	}

	exec := &program{args, nil}
	svc, err := service.New(exec, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err := service.Control(svc, os.Args[1])
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = svc.Run()
	if err != nil {
		log.Fatal(err)
	}

}
