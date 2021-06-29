/*
 * @Author: David
 * @Date: 2020-03-08 12:34:21
 * @LastEditTime: 2020-03-11 18:28:28
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/repositories/key_repository.go
 */

package repositories

import (
	"fmt"
	"marlin/codisclient"
	"marlin/datamodels"
)

// KeyRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type KeyRepository interface {
	GetByID(id, clusterName string) (datamodels.Key, error)
	Set(key datamodels.Key, clusterName string) (datamodels.Key, error)
	DeleteByID(id, clusterName string) (int64, error)
	Expire(id, clusterName string, time int64) (int64, error)
}

// NewKeyRepository returns a new key repository,
// the one and only repository type in our example.
func NewKeyRepository() KeyRepository {
	return &keyRepo{}
}

// keyRepo is a "KeyRepository"
type keyRepo struct {
}

// Select receives a query function
// which is fired for every single user model inside
// our imaginary data source.
// When that function returns true then it stops the iteration.
//
// It returns the query's return last known boolean value
// and the last known user model
// to help callers to reduce the LOC.
//
// It's actually a simple but very clever prototype function
// I'm using everywhere since I firstly think of it,
// hope you'll find it very useful as well.
func (r *keyRepo) GetByID(id, clusterName string) (datamodels.Key, error) {
	value, err := codisclient.Get(clusterName, id)
	if err != nil { //key not exists or other error
		if value == "zk" { //zk error
			return datamodels.Key{
				ID:    id,
				Value: value,
				TTL:   -1,
			}, err
		}
		return datamodels.Key{
			ID:    id,
			Value: value,
			TTL:   0,
		}, err
	}
	if value == "none" { //key no exists
		return datamodels.Key{
			ID:    id,
			Value: value,
			TTL:   0,
		}, err
	}
	ttl, err := codisclient.TTL(clusterName, id)
	return datamodels.Key{
		ID:    id,
		Value: value,
		TTL:   ttl,
	}, err
}

// InsertOrUpdate adds or updates a user to the (memory) storage.
//
// Returns the new user and an error if any.
func (r *keyRepo) Set(key datamodels.Key, clusterName string) (datamodels.Key, error) {
	_, err := codisclient.SetKey(clusterName, key.ID, key.Value)
	if err != nil {
		return datamodels.Key{
			ID:    key.ID,
			Value: key.Value,
			TTL:   0,
		}, err
	}
	ttl, err := codisclient.TTL(clusterName, key.ID)
	return datamodels.Key{
		ID:    key.ID,
		Value: key.Value,
		TTL:   ttl,
	}, err
}

func (r *keyRepo) DeleteByID(id, clusterName string) (int64, error) {
	rtn, err := codisclient.DeleteKey(clusterName, id)
	fmt.Printf("DeleteByID rtn:%d\n", rtn)
	return rtn, err
}

func (r *keyRepo) Expire(id, clusterName string, time int64) (int64, error) {
	rtn, err := codisclient.Expire(clusterName, id, time)
	fmt.Printf("Expire rtn:%d\n", rtn)
	return rtn, err
}
