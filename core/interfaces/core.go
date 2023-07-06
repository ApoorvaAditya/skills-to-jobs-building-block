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

package interfaces

import (
	"application/core/model"

	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
)

// Default exposes client APIs for the driver adapters
type Default interface {
	GetVersion() string
}

// Client exposes client APIs for the driver adapters
type Client interface {
	// OccupationData APIs
	GetOccupationData(code string) (*model.OccupationData, error)
	GetAllOccupationDatas() ([]model.OccupationData, error)

	// UserMatchingResult APIs
	GetUserMatchingResult(id string) (*model.UserMatchingResult, error)
	DeleteUserMatchingResult(id string) error

	// Survey Data APIs
	GetSurveyData(id string) (*model.SurveyData, error)
	CreateSurveyData(surveyData model.SurveyData) (*model.SurveyData, error)
	UpdateSurveyData(surveyData model.SurveyData) error
	DeleteSurveyData(id string) error

	// Occupation Matching
	MatchOccupations(surveyData model.SurveyData, userID string)
}

// Admin exposes administrative APIs for the driver adapters
type Admin interface {
	GetConfig(id string, claims *tokenauth.Claims) (*model.Config, error)
	GetConfigs(configType *string, claims *tokenauth.Claims) ([]model.Config, error)
	CreateConfig(config model.Config, claims *tokenauth.Claims) (*model.Config, error)
	UpdateConfig(config model.Config, claims *tokenauth.Claims) error
	DeleteConfig(id string, claims *tokenauth.Claims) error
}
