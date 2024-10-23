package service

import (
	"context"
	"go.uber.org/zap"
	model "poem_server_admin/internal/model/sys"
	repository "poem_server_admin/internal/repository/sys"
	"poem_server_admin/internal/service"
)

type DictService interface {
	GetDictDataType(ctx context.Context, dictType string) ([]model.SysDictData, error)
}

type dictService struct {
	*service.Service
	repo repository.DictRepository
}

func NewDictService(srv *service.Service, repo repository.DictRepository) DictService {
	return &dictService{
		Service: srv,
		repo:    repo,
	}
}

func (d *dictService) GetDictDataType(ctx context.Context, dictType string) ([]model.SysDictData, error) {
	dataTypeList, err := d.repo.DictDataType(ctx, dictType)
	if err != nil {
		d.Logger.Error(err.Error(), zap.Any("dictType", dictType))
	}
	return dataTypeList, err
}
