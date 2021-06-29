/*
 * @Author: David
 * @Date: 2020-03-08 12:48:18
 * @LastEditTime: 2020-03-12 19:15:07
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/services/key_service.go
 */

package services

import (
	"marlin/datamodels"
	"marlin/repositories"
)

// UserBatchInfoService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserBatchInfoService interface {
	Save(batch datamodels.UserBatchInfo) (result *datamodels.DbResult)
	GetByUserAndFileName(userName, fileName string) (result *datamodels.DbResult)
	GetByUserName(userName string) (result *datamodels.DbResult)
	GetByBatchId(batchId int64) (result *datamodels.DbResult)
	GetUnfinishBatch() (result *datamodels.DbResult)
	UpdateByID(batchInfo datamodels.UserBatchInfo) (result *datamodels.DbResult)
}

// NewUserBatchInfoService returns the default user service.
func NewUserBatchInfoService(repo repositories.UserBatchInfoRepository) UserBatchInfoService {
	return &userBatchInfoService{
		repo: repo,
	}
}

type userBatchInfoService struct {
	repo repositories.UserBatchInfoRepository
}

// GetByID returns a key based on its id.
func (s *userBatchInfoService) GetByUserAndFileName(userName, fileName string) (result *datamodels.DbResult) {
	batch := s.repo.GetByUserAndFileName(userName, fileName)
	result = datamodels.NewDbResult(batch, 0, "get batch info ok.")
	return
}

//GetByUserName ,get user's batch info
func (s *userBatchInfoService) GetByUserName(userName string) (result *datamodels.DbResult) {
	batch := s.repo.GetByUserName(userName)
	result = datamodels.NewDbResult(batch, 0, "get batch info ok.")
	return
}
func (s *userBatchInfoService) GetByBatchId(batchId int64) (result *datamodels.DbResult) {
	batch := s.repo.GetByBatchId(batchId)
	result = datamodels.NewDbResult(batch, 0, "get batch info ok.")
	return
}

//GetUnfinishBatch get all un finished keys
func (s *userBatchInfoService) GetUnfinishBatch() (result *datamodels.DbResult) {
	batchs := s.repo.GetUnfinishBatch()
	result = datamodels.NewDbResult(batchs, 0, "GetUnfinishBatch ok.")
	return
}
func (s *userBatchInfoService) UpdateByID(batchInfo datamodels.UserBatchInfo) (result *datamodels.DbResult) {
	batch := s.repo.UpdateByID(batchInfo)
	result = datamodels.NewDbResult(batch, 0, "GetUnfinishBatch ok.")
	return
}

//Save batch info
func (s *userBatchInfoService) Save(batch datamodels.UserBatchInfo) (result *datamodels.DbResult) {
	app, err := s.repo.Save(batch)
	if err != nil {
		result = datamodels.NewDbResult(app, -1, "Save errpr")
		return
	}
	result = datamodels.NewDbResult(app, 0, "Save ok")
	return

}
