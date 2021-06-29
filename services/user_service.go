/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-11-14 19:39:23
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/services/user_service.go
 */

package services

import (
	"marlin/datamodels"
	"marlin/repositories"
)

// UserService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserService interface {
	GetByNameOrStatus(username string, status int64) (result datamodels.DbResult)
	GetByID(id int64) (result datamodels.DbResult)
	// DeleteByName(username string) (result datamodels.DbResult)
	Save(user datamodels.User) (result datamodels.DbResult, err error)
	UpdateByName(user datamodels.User) (result datamodels.DbResult)
	GetByName(username string) (result datamodels.DbResult)
}

// NewUserService returns the default user service.
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repositories.UserRepository
}

// GetByNameOrStatus returns a user based on its id.
func (s *userService) GetByNameOrStatus(username string, status int64) (result datamodels.DbResult) {
	user := s.repo.GetByNameOrStatus(username, status)
	result.Code = 0
	result.Data = user
	result.Msg = "success"
	return
}

// GetByNameOrStatus returns a user based on its id.
func (s *userService) GetByName(username string) (result datamodels.DbResult) {
	user := s.repo.GetByName(username)
	result.Code = 0
	result.Data = user
	result.Msg = "success"
	return
}

// GetByName returns a user based on its id.
func (s *userService) GetByID(id int64) (result datamodels.DbResult) {
	user := s.repo.GetByID(id)
	result.Code = 0
	result.Data = user
	result.Msg = "success"
	return
}

// Create inserts a new User,
// the userPassword is the client-typed password
// it will be hashed before the insertion to our repository.
func (s *userService) Save(user datamodels.User) (result datamodels.DbResult, err error) {
	user, err = s.repo.Save(user)
	if err != nil {
		result.Code = -1
		result.Msg = "save error."
		return
	}
	result.Code = 0
	result.Data = user
	return
}

func (s *userService) UpdateByName(user datamodels.User) (result datamodels.DbResult) {
	user2 := s.repo.UpdateByName(user)
	result.Code = 0
	result.Data = user2
	result.Msg = "success"
	return
}
