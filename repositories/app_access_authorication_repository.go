/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-11-14 19:38:50
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

// AppAccessAuthoricationRepository handles the basic operations of a user entity/model.
// It's an interface in order to be testable, i.e a memory user repository or
// a connected to an sql database.
type AppAccessAuthoricationRepository interface {
	Save(app datamodels.AppAccessAuthorication) (datamodels.AppAccessAuthorication, error)
	GetByName(applicationName string) (app datamodels.AppAccessAuthorication)
	GetByID(id int64) (app datamodels.AppAccessAuthorication)
	// DeleteByName(username string) (int, datamodels.User)
	GetByApplicant(applicant string) (result *datamodels.DbResult)
	GetAll() (result *datamodels.DbResult)
	UpdateByAppName(app datamodels.AppAccessAuthorication) (datamodels.AppAccessAuthorication, error)
	GetApplications(applicant string, status int64, appName string) (result *datamodels.DbResult)
}

// NewAppAccessAuthoricationRepository returns a new user memory-based repository,
// the one and only repository type in our example.
func NewAppAccessAuthoricationRepository() AppAccessAuthoricationRepository {
	return &appAccessAuthoricationRepository{}
}

// userMemoryRepository is a "UserRepository"
// which manages the users using the memory data source (map).
type appAccessAuthoricationRepository struct {
}

func (repo *appAccessAuthoricationRepository) GetApplications(applicant string, status int64, appName string) (result *datamodels.DbResult) {
	db := datasource.GetDB()
	var apps []datamodels.AppAccessAuthorication
	if applicant != "" {
		if appName != "" {
			db.Where("applicant = ? and application_name= ? and status = ?", applicant, appName, status).Find(&apps)
		} else {
			db.Where("applicant = ? and status = ?", applicant, status).Find(&apps)
		}
	} else {
		if appName != "" {
			db.Where("application_name= ? and status = ?", appName, status).Find(&apps)
		} else {
			db.Where(" status = ?", status).Find(&apps)
		}
	}

	return datamodels.NewDbResult(apps, 0, "get all ok")
}

func (repo *appAccessAuthoricationRepository) GetByName(applicationName string) (app datamodels.AppAccessAuthorication) {
	db := datasource.GetDB()
	db.Where("application_name = ? ", applicationName).Find(&app)
	return
}

//GetByID get user by id
func (repo *appAccessAuthoricationRepository) GetByID(id int64) (app datamodels.AppAccessAuthorication) {
	db := datasource.GetDB()
	db.Where("id = ?", id).Find(&app)
	return
}

// Save to db
func (repo *appAccessAuthoricationRepository) Save(app datamodels.AppAccessAuthorication) (datamodels.AppAccessAuthorication, error) {
	code := 0
	tx := datasource.GetDB().Begin()
	defer utils.Defer(tx, &code)
	var err error
	if err = tx.Save(&app).Error; err != nil {
		log.Println(err)
		code = -1
	}
	return app, err
}
func (repo *appAccessAuthoricationRepository) GetByApplicant(applicant string) (result *datamodels.DbResult) {
	db := datasource.GetDB()
	var apps []datamodels.AppAccessAuthorication
	db.Where("applicant = ? ", applicant).Find(&apps)
	return datamodels.NewDbResult(apps, 0, "get all ok")
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

func (repo *appAccessAuthoricationRepository) GetAll() (result *datamodels.DbResult) {
	db := datasource.GetDB()
	var apps []datamodels.AppAccessAuthorication
	db.Where("1 = 1").Find(&apps)
	return datamodels.NewDbResult(apps, 0, "get all ok")
}

func (repo *appAccessAuthoricationRepository) UpdateByAppName(app datamodels.AppAccessAuthorication) (datamodels.AppAccessAuthorication, error) {
	db := datasource.GetDB()
	db.Model(&app).Where("application_name= ? and status=0", app.ApplicationName).Updates(map[string]interface{}{"status": 1, "secret": app.Secret})
	return app, nil
}
