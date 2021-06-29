/*
 * @Author: your name
 * @Date: 2020-03-08 12:48:18
 * @LastEditTime: 2020-03-12 19:16:22
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/services/key_service.go
 */

package services

import (
	"marlin/datamodels"
	"marlin/repositories"
)

// UserBatchDetailService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserBatchDetailService interface {
	Save(detail datamodels.UserBatchDetail) (result *datamodels.DbResult)
	GetDetailByKeyName(keyName string) (result *datamodels.DbResult)
	GetDetailstByBatchID(batchID int64) (result *datamodels.DbResult)
	UpdateByID(batchInfo datamodels.UserBatchDetail) (result *datamodels.DbResult)
}

// NewUserBatchDetailService returns the default user service.
func NewUserBatchDetailService(repo repositories.UserBatchDetailRepository) UserBatchDetailService {
	return &userBatchDetailService{
		repo: repo,
	}
}

type userBatchDetailService struct {
	repo repositories.UserBatchDetailRepository
}

// GetDetailByKeyName returns a key based on its id.
func (s *userBatchDetailService) GetDetailByKeyName(keyName string) (result *datamodels.DbResult) {
	batch := s.repo.GetByKeyName(keyName)
	result = datamodels.NewDbResult(batch, 0, "get batch info ok.")
	return
}

// GetDetailstByBatchID returns a key based on its id.
func (s *userBatchDetailService) GetDetailstByBatchID(batchID int64) (result *datamodels.DbResult) {
	batch := s.repo.GetByBatchID(batchID)
	result = datamodels.NewDbResult(batch, 0, "get batch info ok.")
	return
}

//Save batch info
func (s *userBatchDetailService) Save(detail datamodels.UserBatchDetail) (result *datamodels.DbResult) {
	app, err := s.repo.Save(detail)
	if err != nil {
		result = datamodels.NewDbResult(app, -1, "Save errpr")
		return
	}
	result = datamodels.NewDbResult(app, 0, "Save ok")
	return

}

//UpdateByID update
func (s *userBatchDetailService) UpdateByID(batchInfo datamodels.UserBatchDetail) (result *datamodels.DbResult) {
	batch := s.repo.UpdateByID(batchInfo)
	result = datamodels.NewDbResult(batch, 0, "UpdateByID ok.")
	return
}
