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

// GetSurveyData finds surveyData by id
func (a Adapter) GetSurveyData(id string) (*model.SurveyData, error) {
	filter := bson.M{"_id": id}

	var data *model.SurveyData
	err := a.db.surveyDatas.FindOne(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeSurveyData, filterArgs(filter), err)
	}

	return data, nil
}

// CreateSurveyData inserts a new surveyData
func (a Adapter) CreateSurveyData(surveyData model.SurveyData) error {
	_, err := a.db.surveyDatas.InsertOne(a.context, surveyData)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionInsert, model.TypeSurveyData, nil, err)
	}

	return nil
}

// UpdateSurveyData updates an surveyData
func (a Adapter) UpdateSurveyData(surveyData model.SurveyData) error {
	filter := bson.M{"_id": surveyData.ID}
	update := bson.M{"$set": bson.M{"scores": surveyData.Scores, "date_updated": time.Now()}}

	_, err := a.db.surveyDatas.UpdateOne(a.context, filter, update, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionUpdate, model.TypeSurveyData, filterArgs(filter), err)
	}
	return nil
}

// DeleteSurveyData deletes an surveyData
func (a Adapter) DeleteSurveyData(id string) error {
	filter := bson.M{"_id": id}

	res, err := a.db.surveyDatas.DeleteOne(a.context, filter, nil)
	if err != nil {
		return errors.WrapErrorAction(logutils.ActionDelete, model.TypeSurveyData, filterArgs(filter), err)
	}
	if res.DeletedCount != 1 {
		return errors.ErrorData(logutils.StatusMissing, model.TypeConfig, filterArgs(filter))
	}

	return nil
}
