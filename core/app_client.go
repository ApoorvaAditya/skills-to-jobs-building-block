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

	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// appClient contains client implementations
type appClient struct {
	app *Application
}

// GetExample gets an Example by ID
func (a appClient) GetExample(orgID string, appID string, id string) (*model.Example, error) {
	return a.app.shared.getExample(orgID, appID, id)
}

// GetOnetData gets an OnetData by ID
func (a appClient) GetOnetDataById(id string) (*model.OnetData, error) {
	return a.app.storage.GetOnetDataById(id)
}

// GetManyOnetData gets Onet data with a specified option
func (a appClient) GetManyOnetData(option *options.FindOptions) (*model.OnetData, error) {
	return a.app.storage.GetManyOnetData(option)
}

// GetBessiData gets an BessiData by ID
func (a appClient) GetBessiData(id string) (*model.BessiData, error) {
	return a.app.storage.GetBessiData(id)
}

// CreateBessiData creates a new BessiData
func (a appClient) CreateBessiData(example model.BessiData) (*model.BessiData, error) {
	example.ID = uuid.NewString()
	err := a.app.storage.CreateBessiData(example)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeBessiData, nil, err)
	}
	return &example, nil
}

// UpdateBessiData updates an BessiData
func (a appClient) UpdateBessiData(example model.BessiData) error {
	return a.app.storage.UpdateBessiData(example)
}

// DeleteBessiData deletes an BessiData by ID
func (a appClient) DeleteBessiData(id string) error {
	return a.app.storage.DeleteBessiData(id)
}

// newAppClient creates new appClient
func newAppClient(app *Application) appClient {
	return appClient{app: app}
}
