package main

import (
	"fmt"
	"os"
	"os/exec"

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
	cfg Config
)

func init() {
	err := configure.Load(app, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(cfg.Jobs) == 0 {
		fmt.Println("Job doesn't defined.")
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
