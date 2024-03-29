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
	"github.com/rokwire/logging-library-go/v2/logutils"
)

const (
	//TypeOccupationData type
	TypeOccupationData logutils.MessageDataType = "occupation data"
	//TypeTechnologySkill type
	TypeTechnologySkill logutils.MessageDataType = "technology skill"
	//TypeWorkstyle type
	TypeWorkstyle logutils.MessageDataType = "workstyle"
)

// OccupationData stores the relevant information about each Occupation from ONET
type OccupationData struct {
	Code             string            `json:"code" bson:"code"`
	Name             string            `json:"name" bson:"name"`
	Description      string            `json:"description" bson:"description"`
	TechnologySkills []TechnologySkill `json:"technology_skills" bson:"technology_skills"`
	Workstyles       []Workstyle       `json:"work_styles" bson:"work_styles"`
}

// TechnologySkill stores the relevant information about each Technology Skill for an occupation
type TechnologySkill struct {
	ID       int      `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name"`
	Examples []string `json:"examples" bson:"examples"`
}

// Workstyle stores the relevant information about each Workstyle for an occupation
type Workstyle struct {
	ID          string  `json:"id" bson:"id"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Scale       string  `json:"scale" bson:"scale"`
	Value       float64 `json:"value" bson:"value"`
}
