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

package core

import (
	"application/core/model"
	"time"

	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// appClient contains client implementations
type appClient struct {
	app *Application
}

// GetExample gets an Example by ID
func (a appClient) GetExample(orgID string, appID string, id string) (*model.Example, error) {
	return a.app.shared.getExample(orgID, appID, id)
}

// GetOccupationData gets an OccupationData by Code
func (a appClient) GetOccupationData(code string) (*model.OccupationData, error) {
	return a.app.storage.GetOccupationData(code)
}

// GetAllOccupationDatas gets all the OccupationDatas
func (a appClient) GetAllOccupationDatas() ([]model.OccupationData, error) {
	return a.app.storage.GetAllOccupationDatas()
}

// GetUserMatchingResult gets an UserMatchingResult by ID
func (a appClient) GetUserMatchingResult(id string) (*model.UserMatchingResult, error) {
	return a.app.storage.GetUserMatchingResult(id)
}

// CreateUserMatchingResult creates a new UserMatchingResult
func (a appClient) CreateUserMatchingResult(userMatchingResult model.UserMatchingResult) (*model.UserMatchingResult, error) {
	userMatchingResult.DateCreated = time.Now()
	err := a.app.storage.CreateUserMatchingResult(userMatchingResult)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeUserMatchingResult, nil, err)
	}
	return &userMatchingResult, nil
}

// UpdateUserMatchingResult updates an UserMatchingResult
func (a appClient) UpdateUserMatchingResult(userMatchingResult model.UserMatchingResult) error {
	return a.app.storage.UpdateUserMatchingResult(userMatchingResult)
}

// DeleteUserMatchingResult deletes an UserMatchingResult by ID
func (a appClient) DeleteUserMatchingResult(id string) error {
	return a.app.storage.DeleteUserMatchingResult(id)
}

// GetSurveyData gets a SurveyData by ID
func (a appClient) GetSurveyData(id string) (*model.SurveyData, error) {
	return a.app.storage.GetSurveyData(id)
}

// CreateSurveyData creates a new SurveyData
func (a appClient) CreateSurveyData(surveyData model.SurveyData) (*model.SurveyData, error) {
	surveyData.ID = uuid.NewString()
	surveyData.DateCreated = time.Now()
	surveyData.Version = "v3.0"
	err := a.app.storage.CreateSurveyData(surveyData)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeSurveyData, nil, err)
	}
	return &surveyData, nil
}

// UpdateSurveyData updates a SurveyData
func (a appClient) UpdateSurveyData(surveyData model.SurveyData) error {
	return a.app.storage.UpdateSurveyData(surveyData)
}

// DeleteSurveyData deletes a SurveyData by ID
func (a appClient) DeleteSurveyData(id string) error {
	return a.app.storage.DeleteSurveyData(id)
}

// GetAllWorkstyleDatas gets all WorkstyleDatas
func (a appClient) GetAllWorkstyleDatas() ([]model.OccupationWorkstyleData, error) {
	return a.app.storage.GetAllWorkstyleDatas()
}

// GetAllWorkstyleDatas finds all WorkstyleDatas for a given occupation
func (a appClient) GetWorkstyleDatasForOccupation(occupationCode string) ([]model.OccupationWorkstyleData, error) {
	return a.app.storage.GetWorkstyleDatasForOccupation(occupationCode)
}

// newAppClient creates new appClient
func newAppClient(app *Application) appClient {
	return appClient{app: app}
}
