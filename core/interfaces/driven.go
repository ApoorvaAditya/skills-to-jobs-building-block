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

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	RegisterStorageListener(listener StorageListener)
	PerformTransaction(func(storage Storage) error) error

	GetConfig(id string) (*model.Config, error)
	SaveConfig(configs model.Config) error
	DeleteConfig(id string) error

	FindExample(orgID string, appID string, id string) (*model.Example, error)
	InsertExample(example model.Example) error
	UpdateExample(example model.Example) error
	DeleteExample(orgID string, appID string, id string) error

	GetBessiData(id string) (*model.BessiData, error)
	CreateBessiData(bessiData model.BessiData) error
	UpdateBessiData(bessiData model.BessiData) error
	DeleteBessiData(id string) error

	GetOnetDataById(id string) (*model.OnetData, error)
	GetManyOnetData(option *options.FindOptions) (*model.OnetData, error)
}

// StorageListener represents storage listener
type StorageListener interface {
	OnConfigsUpdated()
	OnExamplesUpdated()
}
