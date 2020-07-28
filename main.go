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
	Host        string `json:"host"`
	Index       string `json:"index"`
	Interval    int64  `json:"interval"`
	SqlInstance string `json:"sqlinstance"`
}

type program struct {
	args Args
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	p.reqDo()

	ticker := time.NewTicker(time.Second * time.Duration(p.args.Interval))
	for t := range ticker.C {
		_ = t
		p.reqDo()
	}
}

func (p *program) reqDo() {
	dicPerf := ReqCounter(p.args.SqlInstance)
	sb, err := json.Marshal(dicPerf)
	if err != nil {
		log.Fatalf("ReqCounter: %v", err)
		return
	}

	log.Println(string(sb))
	if len(p.args.Host) == 0 {
		return
	}

	cfg := elasticsearch.Config{Addresses: []string{p.args.Host}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
		return
	}

	req := esapi.IndexRequest{
		Index:   p.args.Index + "-" + time.Now().Format("20060102"),
		Body:    strings.NewReader(string(sb)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("reqDo: %v", err)
		return
	}
	defer res.Body.Close()
	log.Println(res)
}

func (p *program) Stop(s service.Service) error {
	return nil
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
		DisplayName: "INTERPARK GoPerf",
	}

	exec := &program{args}
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

	svc.Run()
	err = svc.Run()
	if err != nil {
		log.Fatal(err)
	}

}
