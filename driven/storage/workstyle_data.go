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

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllWorkstyleDatas finds all WorkstyleDatas
func (a Adapter) GetAllWorkstyleDatas() ([]model.OccupationWorkstyleData, error) {
	filter := bson.M{}
	var data []model.OccupationWorkstyleData
	err := a.db.workstyleDatas.Find(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeWorkstyleData, nil, err)
	}

	return data, nil
}

// GetWorkstyleDatasForOccupation finds all WorkstyleDatas
func (a Adapter) GetWorkstyleDatasForOccupation(occupationCode string) ([]model.OccupationWorkstyleData, error) {
	filter := bson.M{"O*NET-SOC Code": occupationCode}
	var data []model.OccupationWorkstyleData
	err := a.db.workstyleDatas.Find(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeWorkstyleData, nil, err)
	}

	return data, nil
}
