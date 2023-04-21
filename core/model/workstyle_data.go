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
	//TypeWorkstyleData example type
	TypeWorkstyleData logutils.MessageDataType = "workstyleData"
)

// WorkstyleData stores statistics such as importance for each workstyle and occupation pair
type OccupationWorkstyleData struct {
	Code        string  `json:"O*NET-SOC Code" bson:"O*NET-SOC Code"`
	Title       string  `json:"Title" bson:"Title"`
	ElementID   string  `json:"Element ID" bson:"Element ID"`
	ElementName string  `json:"Element Name" bson:"Element Name"`
	ScaleID     string  `json:"Scale ID" bson:"Scale ID"`
	ScaleName   string  `json:"Scale Name" bson:"Scale Name"`
	DataValue   float64 `json:"Data Value" bson:"Data Value"`
}
