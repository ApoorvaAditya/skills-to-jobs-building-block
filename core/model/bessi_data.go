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

package model

import (
	"time"

	"github.com/rokwire/logging-library-go/v2/logutils"
)

const (
	//TypeBessiData example type
	TypeBessiData logutils.MessageDataType = "bessi"
)

// BessiData represents the survey results from the BESSI
type BessiData struct {
	ID          string           `json:"id" bson:"_id"`
	Scores      []WorkstyleScore `json:"scores" bson:"scores"`
	DateCreated time.Time        `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time       `json:"date_updated" bson:"date_updated"`
}

// WorkstyleScore represents the score for each workstyle
type WorkstyleScore struct {
	Workstyle string `json:"workstyle" bson:"workstyle"`
	Score     int    `json:"score" bson:"score"`
}
