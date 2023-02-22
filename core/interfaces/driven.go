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
	PerformTransaction(func(torage Storage) error) error

	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error

	GetExample(orgID string, appID string, id string) (*model.Example, error)
	InsertExample(example model.Example) error
	UpdateExample(example model.Example) error
	DeleteExample(orgID string, appID string, id string) error
}

// StorageListener represents storage listener
type StorageListener interface {
	OnConfigsUpdated()
	OnExamplesUpdated()
}
