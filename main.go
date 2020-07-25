package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

type Args struct {
	Host     string `json:"host"`
	Index    string `json:"index"`
	Interval int64  `json:"interval"`
}

type program struct {
	args Args
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	reqDo(p.args.Index, []string{p.args.Host})

	ticker := time.NewTicker(time.Second * time.Duration(p.args.Interval))
	for t := range ticker.C {
		_ = t
		reqDo(p.args.Index, []string{p.args.Host})
	}
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

	log.Println(args)

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
