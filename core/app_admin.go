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
	"application/core/interfaces"
	"application/core/model"
	"time"

	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// appAdmin contains admin implementations
type appAdmin struct {
	app *Application
}

// GetExample gets an Example by ID
func (a appAdmin) GetExample(orgID string, appID string, id string) (*model.Example, error) {
	return a.app.shared.getExample(orgID, appID, id)
}

// CreateExample creates a new Example
func (a appAdmin) CreateExample(example model.Example) (*model.Example, error) {
	example.ID = uuid.NewString()
	err := a.app.storage.InsertExample(example)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeExample, nil, err)
	}
	return &example, nil
}

// UpdateExample updates an Example
func (a appAdmin) UpdateExample(example model.Example) error {
	return a.app.storage.UpdateExample(example)
}

// AppendExample appends to the data in an example - Example of transaction usage
func (a appAdmin) AppendExample(example model.Example) (*model.Example, error) {
	now := time.Now()
	var newExample *model.Example
	transaction := func(storage interfaces.Storage) error {
		oldExample, err := storage.FindExample(example.OrgID, example.AppID, example.ID)
		if err != nil || oldExample == nil {
			return errors.WrapErrorAction(logutils.ActionFind, model.TypeExample, nil, err)
		}

		oldExample.Data = oldExample.Data + "," + example.Data
		oldExample.DateUpdated = &now

		err = storage.UpdateExample(*oldExample)
		if err != nil {
			return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeExample, nil, err)
		}

		newExample = oldExample
		return nil
	}

	err := a.app.storage.PerformTransaction(transaction)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCommit, logutils.TypeTransaction, nil, err)
	}

	return newExample, nil
}

// DeleteExample deletes an Example by ID
func (a appAdmin) DeleteExample(orgID string, appID string, id string) error {
	return a.app.storage.DeleteExample(orgID, appID, id)
}

// newAppAdmin creates new appAdmin
func newAppAdmin(app *Application) appAdmin {
	return appAdmin{app: app}
}
