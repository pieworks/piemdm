package cron

import (
	"log"
	"os"

	"piemdm/internal/service"
	"piemdm/pkg/cron/job"
	client "piemdm/pkg/cron/protocol"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	Scanner       *job.Scanner
	Schedule      *cron.Cron
	cronService   service.CronService
	paramService  service.CronParamService
	entityService service.EntityService
}

func NewCron(scanner *job.Scanner, cronService service.CronService, paramService service.CronParamService, entityService service.EntityService) *Cron {
	return &Cron{
		Scanner: scanner,
		// SkipIfStillRunning skips an invocation of the Job if a previous invocation is still running. It logs skips to the given logger at Info level.
		Schedule:      cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))),
		cronService:   cronService,
		paramService:  paramService,
		entityService: entityService,
	}
}

func (s *Cron) Start() error {
	where := map[string]any{}
	where["status"] = "Normal"
	jobs, err1 := s.cronService.Find("", where)
	if err1 != nil {
		return err1
	}
	for _, job := range jobs {
		paramWhere := map[string]any{}
		paramWhere["status"] = "Normal"
		paramWhere["cron_code"] = job.Code
		params, err2 := s.paramService.Find("", paramWhere)
		if err2 != nil {
			return err2
		}

		_, err := s.Schedule.AddJob(job.Expression, client.NewProtocol(job.Protocol, job, params, s.entityService))
		if err != nil {
			return err
		}
	}
	// _, err := s.Schedule.AddJob("0 */30 * * * *", s.Scanner)
	s.Schedule.Start()

	// print a snapshot of the cron entries.
	// fmt.Printf("\ns.Schedule.Entries(): %#v\n\n", s.Schedule.Entries())
	return nil
}
