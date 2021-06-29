/*
 * @Author: David
 * @Date: 2020-03-09 15:18:48
 * @LastEditTime: 2020-03-12 00:02:12
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/services/app_access_authorication_service.go
 */
package services

import (
	"marlin/datamodels"
	"marlin/repositories"
)

//AppAccessAuthoricationService application access service
type AppAccessAuthoricationService interface {
	Save(appAccessAuthorication datamodels.AppAccessAuthorication) (result *datamodels.DbResult)
	GetByName(appName string) (result *datamodels.DbResult)
	GetByID(id int64) (result *datamodels.DbResult)
	GetByApplicant(applicant string) (result *datamodels.DbResult)
	GetAll() (result *datamodels.DbResult)
	UpdateByAppName(app datamodels.AppAccessAuthorication) (result *datamodels.DbResult)
	GetApplications(applicant string, status int64, appName string) (result *datamodels.DbResult)
}

type appAccessAuthoricationServices struct {
	repo repositories.AppAccessAuthoricationRepository
}

//NewAppAccessAuthoricationService new service
func NewAppAccessAuthoricationService(repo repositories.AppAccessAuthoricationRepository) AppAccessAuthoricationService {
	return &appAccessAuthoricationServices{
		repo: repo,
	}
}
func (u *appAccessAuthoricationServices) GetApplications(applicant string, status int64, appName string) (result *datamodels.DbResult) {
	application := u.repo.GetApplications(applicant, status, appName)
	// result = datamodels.NewDbResult(application, 0, "GetApplications ok")
	return application
}

func (u *appAccessAuthoricationServices) GetByName(appName string) (result *datamodels.DbResult) {
	application := u.repo.GetByName(appName)
	result = datamodels.NewDbResult(application, 0, "GetApplications ok")
	return
}

func (u appAccessAuthoricationServices) Save(appAccessAuthorication datamodels.AppAccessAuthorication) (result *datamodels.DbResult) {
	app, err := u.repo.Save(appAccessAuthorication)
	if err != nil {
		result = datamodels.NewDbResult(app, -1, "Save errpr")
		return
	}
	result = datamodels.NewDbResult(app, 0, "Save ok")
	return
}

func (u *appAccessAuthoricationServices) GetByID(id int64) (result *datamodels.DbResult) {
	app := u.repo.GetByID(id)
	result = datamodels.NewDbResult(app, 0, "GetApplications ok")
	return
}

func (u *appAccessAuthoricationServices) GetByApplicant(applicant string) (result *datamodels.DbResult) {
	return u.repo.GetByApplicant(applicant)
}

func (u *appAccessAuthoricationServices) GetAll() (result *datamodels.DbResult) {
	return u.repo.GetAll()
}
func (u *appAccessAuthoricationServices) UpdateByAppName(app datamodels.AppAccessAuthorication) (result *datamodels.DbResult) {
	app, err := u.repo.UpdateByAppName(app)
	if err != nil {
		result = datamodels.NewDbResult(app, -1, "Save errpr")
		return
	}
	result = datamodels.NewDbResult(app, 0, "Save ok")
	return
}
