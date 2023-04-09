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
	//TypeUserMatchingResult type
	TypeUserMatchingResult logutils.MessageDataType = "userMatchingResult"
	//TypeMatch type
	TypeMatch logutils.MessageDataType = "match"
)

// UserMatchingResult represents the matching results of a specific user
type UserMatchingResult struct {
	ID          string     `json:"id" bson:"_id"`
	Version     string     `json:"version" bson:"version"`
	Matches     []Match    `json:"matches" bson:"matches"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
}

// Match represents a occupation match and the corresponding score
type Match struct {
	OccupationCode string         `json:"occupation_code" bson:"occupation_code"`
	Occupation     OccupationData `json:"occupation" bson:"occupation"`
	Score          float64        `json:"score" bson:"score"`
}
