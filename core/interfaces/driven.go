// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interfaces

import (
	"application/core/model"
)

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	RegisterStorageListener(listener StorageListener)
	PerformTransaction(func(storage Storage) error) error

	FindConfig(configType string, appID string, orgID string) (*model.Config, error)
	FindConfigByID(id string) (*model.Config, error)
	FindConfigs(configType *string) ([]model.Config, error)
	InsertConfig(config model.Config) error
	UpdateConfig(config model.Config) error
	DeleteConfig(id string) error

	GetOccupationData(id string) (*model.OccupationData, error)
	GetAllOccupationDatas() ([]model.OccupationData, error)

	GetUserMatchingResult(id string) (*model.UserMatchingResult, error)
	SaveUserMatchingResult(bessiData model.UserMatchingResult) error
	DeleteUserMatchingResult(id string) error

	GetSurveyData(id string) (*model.SurveyData, error)
	CreateSurveyData(surveyData model.SurveyData) error
	UpdateSurveyData(surveyData model.SurveyData) error
	DeleteSurveyData(id string) error
}

// StorageListener represents storage listener
type StorageListener interface {
	OnConfigsUpdated()
}
