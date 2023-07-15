package cron

import (
	"context"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	cron *cron.Cron
	jobs []Job
}

type Job struct {
	job  func(ctx context.Context) error
	spec string
}

func New() *Cron {
	return &Cron{
		cron: cron.New(),
	}
}

// `func (c *Cron) AddJob(spec string, job func(ctx context.Context) error)` is a method of the `Cron`
// struct that adds a new job to the list of jobs to be executed by the cron scheduler. It takes two
// parameters: `spec` which is a string representing the cron schedule for the job, and `job` which is
// a function that takes a context as a parameter and returns an error. The method first gets the name
// of the function using the `runtime.FuncForPC` and `reflect.ValueOf` functions. It then parses the
// cron schedule using the `cron.ParseStandard` function and adds the job to the list of jobs. Finally,
// it logs a message indicating that the job has been registered and when it will be executed next.
func (c *Cron) AddJob(spec string, job func(ctx context.Context) error) {
	name := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()
	schedule, err := cron.ParseStandard(spec)
	if err != nil {
		log.Printf("Failed to register job %s becasue of invalid spec %s %s", name, spec, err.Error())
		return
	}
	c.jobs = append(c.jobs, Job{job: job, spec: spec})
	log.Printf("Job %s is registered. Next execution at %s", name, schedule.Next(time.Now()))
}

// `func (c *Cron) Start(ctx context.Context)` is a method of the `Cron` struct that starts the cron
// scheduler and schedules all the jobs added to the list of jobs to be executed. It takes a context as
// a parameter and uses it to execute the jobs. The method iterates over the list of jobs and adds each
// job to the cron scheduler using the `AddFunc` method of the `cron` package. Finally, it starts the
// cron scheduler using the `Start` method of the `cron` package. Once started, the cron scheduler will
// execute the jobs at their scheduled times according to their respective cron schedules.
func (c *Cron) Start(ctx context.Context) {
	for _, job := range c.jobs {
		id, err := c.cron.AddFunc(job.spec, func() {
			job.job(ctx)
		})
		if err != nil {
			log.Printf("Failed to AddFunc to cron for id %+v because %s", id, err.Error())
		}
	}
	c.cron.Start()
}
