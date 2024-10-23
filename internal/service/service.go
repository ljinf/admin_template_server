package service

import (
	"poem_server_admin/internal/repository"
	"poem_server_admin/pkg/jwt"
	"poem_server_admin/pkg/log"
	"poem_server_admin/pkg/sid"
)

type Service struct {
	Logger *log.Logger
	Sid    *sid.Sid
	Jwt    *jwt.JWT
	Tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		Logger: logger,
		Sid:    sid,
		Jwt:    jwt,
		Tm:     tm,
	}
}
