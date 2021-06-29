/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-03-12 19:13:41
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/repositories/user_repository.go
 */
package repositories

import (
	"log"
	"marlin/datamodels"
	"marlin/datasource"
	"marlin/utils"
)

// UserBatchInfoRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type UserBatchInfoRepository interface {
	Save(batch datamodels.UserBatchInfo) (datamodels.UserBatchInfo, error)
	GetByUserAndFileName(userName, fileName string) (app datamodels.UserBatchInfo)
	GetByUserName(userName string) (batchs []datamodels.UserBatchInfo)
	GetByBatchId(batchId int64) (batchs []datamodels.UserBatchInfo)
	GetUnfinishBatch() (batchs []datamodels.UserBatchInfo)
	UpdateByID(batch datamodels.UserBatchInfo) datamodels.UserBatchInfo
}

// NewUserBatchInfoRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserBatchInfoRepository() UserBatchInfoRepository {
	return &userBatchInfoRepository{}
}

// userBatchInfoRepository is a "UserBatchInfoRepository"
// which manages the users using the memory data source (map).
type userBatchInfoRepository struct {
}

func (repo *userBatchInfoRepository) GetByUserAndFileName(userName, fileName string) (app datamodels.UserBatchInfo) {
	db := datasource.GetDB()
	db.Where("user_name = ? and file_name = ? and status = ?", userName, fileName, 0).Limit(1).Find(&app)
	return
}
func (repo *userBatchInfoRepository) GetByUserName(userName string) (batchs []datamodels.UserBatchInfo) {
	db := datasource.GetDB()
	db.Where("user_name = ?", userName).Find(&batchs)
	return
}
func (repo *userBatchInfoRepository) GetByBatchId(batchId int64) (batchs []datamodels.UserBatchInfo) {
	db := datasource.GetDB()
	db.Where("id = ?", batchId).Find(&batchs)
	return
}
func (repo *userBatchInfoRepository) GetUnfinishBatch() (batchs []datamodels.UserBatchInfo) {
	db := datasource.GetDB()
	db.Where("status = ?", 0).Find(&batchs)
	return
}

//UpdateByID finished
func (repo *userBatchInfoRepository) UpdateByID(batch datamodels.UserBatchInfo) datamodels.UserBatchInfo {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)

	if err := tx.Model(&batch).Where("id = ?", batch.ID).UpdateColumn("status", 1).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return batch
}

// Save to db
func (repo *userBatchInfoRepository) Save(batch datamodels.UserBatchInfo) (datamodels.UserBatchInfo, error) {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)
	var err error
	if err = tx.Save(&batch).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return batch, err
}
