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
	//TypeOccupationData example type
	TypeOccupationData logutils.MessageDataType = "occupation"
)

// OccupationData represents the survey results from the Occupation
type OccupationData struct {
	Code       string       `json:"code" bson:"code"`		
	Title      string 		`json:"title" bson:"title"`
	Description string		`json:"description" bson:"description"`
}