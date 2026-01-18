//go:build wireinject
// +build wireinject

package main

import (
	"piemdm/internal/repository"
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/log"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewBaseRepository,
	repository.NewUserRepository,
)

func newApp(*viper.Viper, *log.Logger) (*Migrate, func(), error) {
	panic(wire.Build(
		RepositorySet,
		NewMigrate,
		sid.NewSid,
	))
}
