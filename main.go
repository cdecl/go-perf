package main

import (
	"flag"
	"log"
	"time"

	"github.com/kardianos/service"
)

type Args struct {
	Index    string
	Host     string
	Interval int64
}

func usage() (Args, bool) {
	args := Args{}
	flag.StringVar(&args.Index, "i", "", "Elasticsearch index name (require)")
	flag.StringVar(&args.Host, "h", "", "Elasticsearch host (require) e.g. http://localhost:9200")
	flag.Int64Var(&args.Interval, "t", 20, "Interval : default :20 (seconds)")
	flag.Parse()

	isFlagPassed := func(name string) bool {
		found := false
		flag.Visit(func(f *flag.Flag) {
			if f.Name == name {
				found = true
			}
		})
		return found
	}

	found := isFlagPassed("i")
	found = found && isFlagPassed("h")

	if !found {
		flag.Usage()
		return args, false
	}
	return args, true
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

func main() {
	args, ok := usage()
	if !ok {
		return
	}

	svcConfig := &service.Config{
		Name:        "GoPerf",
		DisplayName: "GoPerf",
	}

	exec := &program{args}
	svc, err := service.New(exec, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = svc.Run()

	if err != nil {
		logger.Error(err)
	}

}
