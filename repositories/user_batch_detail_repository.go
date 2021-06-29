/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-11-14 19:38:59
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

// UserBatchDetailRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type UserBatchDetailRepository interface {
	Save(batch datamodels.UserBatchDetail) (datamodels.UserBatchDetail, error)
	GetByKeyName(keyName string) (app datamodels.UserBatchDetail)
	GetByBatchID(batchID int64) (details []datamodels.UserBatchDetail)
	UpdateByID(batch datamodels.UserBatchDetail) datamodels.UserBatchDetail
}

// NewUserBatchDetailRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserBatchDetailRepository() UserBatchDetailRepository {
	return &userBatchDetailRepository{}
}

// userBatchDetailRepository is a "UserBatchDetailRepository"
// which manages the users using the memory data source (map).
type userBatchDetailRepository struct {
}

func (repo *userBatchDetailRepository) GetByKeyName(keyName string) (app datamodels.UserBatchDetail) {
	db := datasource.GetDB()
	db.Where("key_name = ?", keyName).Find(&app)
	return
}

// Save to db
func (repo *userBatchDetailRepository) Save(batch datamodels.UserBatchDetail) (datamodels.UserBatchDetail, error) {
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

//GetByBatchID get all details
func (repo *userBatchDetailRepository) GetByBatchID(batchID int64) (details []datamodels.UserBatchDetail) {
	db := datasource.GetDB()
	db.Where("batch_id = ? ", batchID).Find(&details)
	return
}

//UpdateByID update rec
func (repo *userBatchDetailRepository) UpdateByID(batch datamodels.UserBatchDetail) datamodels.UserBatchDetail {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)

	if err := tx.Model(&batch).Where("id = ?", batch.ID).UpdateColumn("status", 1).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return batch
}
