// Package service implements the business logic layer for the PieMDM system.
// It contains service interfaces and implementations that handle business rules,
// data validation, and coordinate between the handler and repository layers.
package service

import (
	"piemdm/pkg/helper/sid"
	"piemdm/pkg/jwt"
	"piemdm/pkg/log"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
}

func NewService(logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
	}
}
