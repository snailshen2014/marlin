/*
 * @Author: David
 * @Date: 2020-03-08 12:48:18
 * @LastEditTime: 2020-11-14 19:39:16
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/services/key_service.go
 */

package services

import (
	"marlin/datamodels"
	"marlin/repositories"
)

// KeyService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type KeyService interface {
	GetByID(id, clusterName string) (datamodels.Key, error)
	DeleteByID(id, clusterName string) (int64, error)
	Set(key datamodels.Key, clusterName string) (datamodels.Key, error)
	Expire(id, clusterName string, time int64) (int64, error)
}

// NewKeyService returns the default user service.
func NewKeyService(repo repositories.KeyRepository) KeyService {
	return &keyService{
		repo: repo,
	}
}

type keyService struct {
	repo repositories.KeyRepository
}

// GetByID returns a key based on its id.
func (s *keyService) GetByID(id, clusterName string) (datamodels.Key, error) {
	return s.repo.GetByID(id, clusterName)
}

func (s *keyService) Set(key datamodels.Key, clusterName string) (datamodels.Key, error) {
	return s.repo.Set(key, clusterName)
}

// DeleteByID deletes a key by its id.
//
// Returns true if deleted otherwise false.
func (s *keyService) DeleteByID(id, clusterName string) (int64, error) {
	return s.repo.DeleteByID(id, clusterName)
}

func (s *keyService) Expire(id, clusterName string, time int64) (int64, error) {
	return s.repo.Expire(id, clusterName, time)
}
