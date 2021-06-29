/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-11-14 19:39:06
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

// UserRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type UserRepository interface {
	Save(user datamodels.User) (datamodels.User, error)
	GetByNameOrStatus(username string, status int64) (users []datamodels.User)
	GetByName(username string) (user datamodels.User)
	GetByID(id int64) (user datamodels.User)
	// DeleteByName(username string) (int, datamodels.User)
	UpdateByName(user datamodels.User) datamodels.User
}

// NewUserRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// userMemoryRepository is a "UserRepository"
// which manages the users using the memory data source (map).
type userRepository struct {
}

func (repo *userRepository) GetByNameOrStatus(username string, status int64) (users []datamodels.User) {
	db := datasource.GetDB()
	if username == "" {
		db.Where("status = ?", status).Find(&users)
		return
	}
	db.Where("user_name = ? and status = ?", username, status).Find(&users)
	return
}
func (repo *userRepository) GetByName(username string) (user datamodels.User) {
	db := datasource.GetDB()

	db.Where("user_name = ? ", username).Find(&user)
	return
}

//GetByID get user by id
func (repo *userRepository) GetByID(id int64) (user datamodels.User) {
	db := datasource.GetDB()
	db.Where("id = ?", id).Find(&user)
	return
}

// Save to db
func (repo *userRepository) Save(user datamodels.User) (datamodels.User, error) {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)
	var err error
	if err = tx.Save(&user).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return user, err
}

// func (repo *userRepository) DeleteByName(username string) (int, datamodels.User) {
// 	code := 0
// 	tx := datasource.GetDB().Begin()
// 	defer utils.Defer(tx, &code)

// 	if err := tx.Save(&user).Error; err != nil {
// 		log.Println(err)
// 		code = -1
// 	}
// 	return code, user
// }
func (repo *userRepository) UpdateByName(user datamodels.User) datamodels.User {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)

	if err := tx.Model(&user).Where("user_name= ?", user.UserName).UpdateColumn("status", 1).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return user
}
