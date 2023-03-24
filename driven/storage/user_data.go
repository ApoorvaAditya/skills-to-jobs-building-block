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

package storage

import (
	"application/core/model"
	"time"

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetUserData finds userData by id
func (a Adapter) GetUserData(id string) (*model.UserData, error) {
	filter := bson.M{"_id": id}

	var data *model.UserData
	err := a.db.userDatas.FindOneWithContext(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeUserData, filterArgs(filter), err)
	}

	return data, nil
}

// CreateUserData inserts a new userData
func (a Adapter) CreateUserData(userData model.UserData) error {
	_, err := a.db.userDatas.InsertOneWithContext(a.context, userData)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInsert, model.TypeUserData, nil, err)
	}

	return nil
}

// UpdateUserData updates an userData
func (a Adapter) UpdateUserData(userData model.UserData) error {
	filter := bson.M{"_id": userData.ID}
	update := bson.M{"$set": bson.M{"matches": userData.Matches, "data_updated": time.Now()}}

	_, err := a.db.userDatas.UpdateOneWithContext(a.context, filter, update, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeUserData, filterArgs(filter), err)
	}
	return nil
}

// DeleteUserData deletes an userData
func (a Adapter) DeleteUserData(id string) error {
	filter := bson.M{"_id": id}

	res, err := a.db.userDatas.DeleteOneWithContext(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeUserData, filterArgs(filter), err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, filterArgs(filter))
	}

	return nil
}
