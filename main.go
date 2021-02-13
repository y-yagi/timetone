package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/0xAX/notificator"
	"github.com/robfig/cron"
	"github.com/y-yagi/configure"
)

type Job struct {
	Spec    string
	Message string
	Command []string
}

type Config struct {
	Jobs []Job
}

const (
	app = "timetone"
)

var (
	cfg   Config
	flags *flag.FlagSet
	list  bool
)

func init() {
	err := configure.Load(app, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setFlags() {
	flags = flag.NewFlagSet(app, flag.ExitOnError)
	flags.BoolVar(&list, "l", false, "list as crontab format")
}

func main() {
	setFlags()
	flags.Parse(os.Args[1:])
	if len(cfg.Jobs) == 0 {
		fmt.Println("Job doesn't defined.")
		return
	}

	if list {
		printList()
		return
	}

	done := make(chan bool)
	notify := notificator.New(notificator.Options{})
	var err error

	c := cron.New()
	for _, job := range cfg.Jobs {
		err = c.AddFunc(job.Spec, func() {
			msg := job.Message
			if len(msg) == 0 {
				msg = fmt.Sprintf("Run %v", job.Command)
			}
			notify.Push("TimeTone", msg, "", notificator.UR_NORMAL)

			if len(job.Command) != 0 {
				cmd := exec.Command(job.Command[0], job.Command[1:]...)
				cmd.Run()
			}
		})

		if err != nil {
			fmt.Printf("Job setting failed: %v\n", err)
			return
		}
	}

	c.Start()

	<-done
}

func printList() {
	buf := ""
	for _, job := range cfg.Jobs {
		buf += fmt.Sprintf("%s %s\n", job.Spec, strings.Join(job.Command, " "))
	}
	fmt.Printf("%s", buf)
}
