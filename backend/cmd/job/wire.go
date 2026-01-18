//go:build wireinject
// +build wireinject

package main

import (
	"piemdm/internal/repository"
	"piemdm/internal/service"
	"piemdm/pkg/cron/job"
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/jwt"
	"piemdm/pkg/log"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewCronService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewBaseRepository,
	repository.NewUserRepository,
	repository.NewUserRoleRepository,
	repository.NewCronRepository,
)

var JobSet = wire.NewSet(job.NewScanner)

func newApp(*viper.Viper, *log.Logger) (*job.Scanner, func(), error) {
	panic(wire.Build(
		RepositorySet,
		ServiceSet,
		JobSet,
		sid.NewSid,
		jwt.NewJwt,
	))
}
