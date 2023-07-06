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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetUserMatchingResult finds userMatchingResult by id
func (a Adapter) GetUserMatchingResult(id string) (*model.UserMatchingResult, error) {
	filter := bson.M{"_id": id}

	var data *model.UserMatchingResult
	err := a.db.matchResults.FindOne(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeUserMatchingResult, filterArgs(filter), err)
	}

	return data, nil
}

// SaveUserMatchingResult saves a userMatchingResult
func (a Adapter) SaveUserMatchingResult(userMatchingResult model.UserMatchingResult) error {
	filter := bson.M{"_id": userMatchingResult.ID}
	update := bson.M{
		"$set": bson.M{
			"matches":      userMatchingResult.Matches,
			"date_updated": time.Now().UTC(),
		},
		"$setOnInsert": bson.M{
			"date_created": time.Now().UTC(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := a.db.matchResults.UpdateOne(a.context, filter, update, opts)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeUserMatchingResult, filterArgs(filter), err)
	}
	return nil
}

// DeleteUserMatchingResult deletes an userMatchingResult
func (a Adapter) DeleteUserMatchingResult(id string) error {
	filter := bson.M{"_id": id}

	res, err := a.db.matchResults.DeleteOne(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeUserMatchingResult, filterArgs(filter), err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, filterArgs(filter))
	}

	return nil
}
