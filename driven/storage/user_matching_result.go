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

// GetUserMatchingResult finds userMatchingResult by id
func (a Adapter) GetUserMatchingResult(id string) (*model.UserMatchingResult, error) {
	filter := bson.M{"_id": id}

	var data *model.UserMatchingResult
	err := a.db.userMatchingResults.FindOne(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeUserMatchingResult, filterArgs(filter), err)
	}

	return data, nil
}

// CreateUserMatchingResult inserts a new userMatchingResult
func (a Adapter) CreateUserMatchingResult(userMatchingResult model.UserMatchingResult) error {
	_, err := a.db.userMatchingResults.InsertOne(a.context, userMatchingResult)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInsert, model.TypeUserMatchingResult, nil, err)
	}

	return nil
}

// UpdateUserMatchingResult updates an userMatchingResult
func (a Adapter) UpdateUserMatchingResult(userMatchingResult model.UserMatchingResult) error {
	filter := bson.M{"_id": userMatchingResult.ID}
	update := bson.M{"$set": bson.M{"matches": userMatchingResult.Matches, "data_updated": time.Now()}}

	_, err := a.db.userMatchingResults.UpdateOne(a.context, filter, update, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeUserMatchingResult, filterArgs(filter), err)
	}
	return nil
}

// DeleteUserMatchingResult deletes an userMatchingResult
func (a Adapter) DeleteUserMatchingResult(id string) error {
	filter := bson.M{"_id": id}

	res, err := a.db.userMatchingResults.DeleteOne(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeUserMatchingResult, filterArgs(filter), err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, filterArgs(filter))
	}

	return nil
}
