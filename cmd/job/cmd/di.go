package cmd

import (
	esrepo "github.com/GodYao1995/Goooooo/internal/admin/es"
	"github.com/GodYao1995/Goooooo/internal/admin/repository"
	job "github.com/GodYao1995/Goooooo/internal/job"
	tasks "github.com/GodYao1995/Goooooo/internal/job/tasks"
	"github.com/GodYao1995/Goooooo/pkg/db"
	"github.com/GodYao1995/Goooooo/pkg/logger"
	"github.com/GodYao1995/Goooooo/pkg/xes"
	xjob "github.com/GodYao1995/Goooooo/pkg/xjob"
	"github.com/RichardKnop/machinery/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var JobWorkerModule = fx.Invoke(JobWorker)

func Run() {
	fx.New(inject()).Run()
}

func JobWorker(jobServer *machinery.Server, _jobs *tasks.UserESTask) {
	jobServer.RegisterTask("sum", tasks.Sum)
	jobServer.RegisterTask("user2es", _jobs.UsersToES)
	worker := jobServer.NewWorker("worker", 0)
	worker.Launch()
}

func inject() fx.Option {
	return fx.Options(
		// Provide
		configModule,
		logger.Module,
		db.Module,
		xjob.Module,
		xes.Module,
		repository.Module,
		esrepo.Module,
		job.Module,
		// Invoke
		JobWorkerModule,
		// Options
		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)
}
