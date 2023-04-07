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

// GetUserData gets an UserData by ID
func (a appClient) GetUserData(id string) (*model.UserData, error) {
	return a.app.storage.GetUserData(id)
}

// CreateUserData creates a new UserData
func (a appClient) CreateUserData(userData model.UserData) (*model.UserData, error) {
	userData.ID = uuid.NewString()
	userData.DateCreated = time.Now()
	err := a.app.storage.CreateUserData(userData)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeUserData, nil, err)
	}
	return &userData, nil
}

// UpdateUserData updates an UserData
func (a appClient) UpdateUserData(userData model.UserData) error {
	return a.app.storage.UpdateUserData(userData)
}

// DeleteUserData deletes an UserData by ID
func (a appClient) DeleteUserData(id string) error {
	return a.app.storage.DeleteUserData(id)
}

// newAppClient creates new appClient
func newAppClient(app *Application) appClient {
	return appClient{app: app}
}
