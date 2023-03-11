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

// GetBessiData finds bessiData by id
func (a Adapter) GetBessiData(id string) (*model.BessiData, error) {
	filter := bson.M{"_id": id}

	var data *model.BessiData
	err := a.db.bessiDatas.FindOneWithContext(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeBessiData, filterArgs(filter), err)
	}

	return data, nil
}

// CreateBessiData inserts a new bessiData
func (a Adapter) CreateBessiData(bessiData model.BessiData) error {
	_, err := a.db.bessiDatas.InsertOneWithContext(a.context, bessiData)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInsert, model.TypeBessiData, nil, err)
	}

	return nil
}

// UpdateBessiData updates an bessiData
func (a Adapter) UpdateBessiData(bessiData model.BessiData) error {
	filter := bson.M{"_id": bessiData.ID}
	update := bson.M{"$set": bson.M{"data": bessiData.Data, "date_updated": time.Now()}}

	_, err := a.db.bessiDatas.UpdateOneWithContext(a.context, filter, update, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeBessiData, filterArgs(filter), err)
	}
	return nil
}

// DeleteBessiData deletes an bessiData
func (a Adapter) DeleteBessiData(id string) error {
	filter := bson.M{"_id": id}

	res, err := a.db.bessiDatas.DeleteOneWithContext(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeBessiData, filterArgs(filter), err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, filterArgs(filter))
	}

	return nil
}
