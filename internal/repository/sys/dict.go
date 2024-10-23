package sys

import (
	"context"
	model "poem_server_admin/internal/model/sys"
	"poem_server_admin/internal/repository"
)

type DictRepository interface {
	DictDataList(ctx context.Context) ([]model.SysDictData, error)
	DictDataType(ctx context.Context, dictType string) ([]model.SysDictData, error)
}

type dictRepository struct {
	*repository.Repository
}

func NewDictRepository(repo *repository.Repository) DictRepository {
	return &dictRepository{
		Repository: repo,
	}
}

func (d *dictRepository) DictDataList(ctx context.Context) ([]model.SysDictData, error) {
	return nil, nil
}

func (d *dictRepository) DictDataType(ctx context.Context, dictType string) ([]model.SysDictData, error) {
	var list []model.SysDictData
	err := d.DB(ctx).Where("status = ? and dict_type = ?", 0, dictType).Order("dict_sort asc").Find(&list).Error
	return list, err
}
