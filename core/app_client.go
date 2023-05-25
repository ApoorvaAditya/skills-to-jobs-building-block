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

package core

import (
	"application/core/model"
	"sort"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/google/uuid"
	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// appClient contains client implementations
type appClient struct {
	app *Application
}

// GetOccupationData gets an OccupationData by Code
func (a appClient) GetOccupationData(code string) (*model.OccupationData, error) {
	return a.app.storage.GetOccupationData(code)
}

// GetAllOccupationDatas gets all the OccupationDatas
func (a appClient) GetAllOccupationDatas() ([]model.OccupationData, error) {
	return a.app.storage.GetAllOccupationDatas()
}

// GetUserMatchingResult gets an UserMatchingResult by ID
func (a appClient) GetUserMatchingResult(id string) (*model.UserMatchingResult, error) {
	return a.app.storage.GetUserMatchingResult(id)
}

// DeleteUserMatchingResult deletes an UserMatchingResult by ID
func (a appClient) DeleteUserMatchingResult(id string) error {
	return a.app.storage.DeleteUserMatchingResult(id)
}

// GetSurveyData gets a SurveyData by ID
func (a appClient) GetSurveyData(id string) (*model.SurveyData, error) {
	return a.app.storage.GetSurveyData(id)
}

// CreateSurveyData creates a new SurveyData
func (a appClient) CreateSurveyData(surveyData model.SurveyData) (*model.SurveyData, error) {
	surveyData.ID = uuid.NewString()
	surveyData.DateCreated = time.Now()
	surveyData.Version = "v3.0"
	err := a.app.storage.CreateSurveyData(surveyData)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionCreate, model.TypeSurveyData, nil, err)
	}
	return &surveyData, nil
}

// UpdateSurveyData updates a SurveyData
func (a appClient) UpdateSurveyData(surveyData model.SurveyData) error {
	return a.app.storage.UpdateSurveyData(surveyData)
}

// DeleteSurveyData deletes a SurveyData by ID
func (a appClient) DeleteSurveyData(id string) error {
	return a.app.storage.DeleteSurveyData(id)
}

func (a appClient) MatchOccupations(surveyData model.SurveyData, userID string) {
	occupations, err := a.GetAllOccupationDatas()
	if err != nil {
		return
	}

	matches := a.runMatchingAlgo(surveyData.Scores, occupations)
	userMatchingResult := model.UserMatchingResult{
		ID:      userID,
		Matches: matches,
		Version: surveyData.Version,
	}

	a.app.storage.SaveUserMatchingResult(userMatchingResult)
}

func (a appClient) runMatchingAlgo(userScores []model.WorkstyleScore, occupations []model.OccupationData) []model.Match {
	matches := make([]model.Match, 0)
	for _, occupation := range occupations {
		if len(occupation.Workstyles) > 0 {
			occupationMatch := a.runMatchingAlgorithmPerOccupation(occupation, userScores)
			matches = append(matches, occupationMatch)
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].MatchPercent > matches[j].MatchPercent
	})
	return matches
}

// It runs the matching alogrithm on the entire occupation list and its workstyles to find the p-value for each occupation
func (a appClient) runMatchingAlgorithmPerOccupation(occupation model.OccupationData, userScores []model.WorkstyleScore) model.Match {
	match := model.Match{Occupation: model.OccupationMatch{Code: occupation.Code, Name: occupation.Name}}

	dfUserScoresUnsorted := dataframe.LoadStructs(userScores)
	dfUserScores := dfUserScoresUnsorted.Arrange(dataframe.Sort("Score"))
	dfImportanceUnsorted := dataframe.LoadStructs(occupation.Workstyles)
	dfImportance := dfImportanceUnsorted.Arrange(dataframe.Sort("Value"))

	bessiToWorkstyles := map[string]string{
		"stress_regulation":         "Stress Tolerance",
		"adaptability":              "Adaptability/Flexibility",
		"capacity_social_warmth":    "Concern for Others",
		"abstract_thinking":         "Analytical Thinking",
		"teamwork":                  "Cooperation",
		"responsibility_management": "Dependability",
		"detail_management":         "Attention to Detail",
		"initiative":                "Initiative",
		"anger_management":          "Self-Control",
		"capacity_consistency":      "Persistence",
		"capacity_independence":     "Independence",
		"perspective_taking":        "Social Orientation",
		"goal_regulation":           "Achievement/Effort",
		"creativity":                "Innovation",
		"ethical_competence":        "Integrity",
		"leadership":                "Leadership",
	}

	sumSquared := 0.0
	n := float64(dfUserScores.Nrow())
	for i := 0; i < dfUserScores.Nrow(); i++ {
		row := dfUserScores.Subset(i)
		workstyle := bessiToWorkstyles[row.Col("Workstyle").Elem(0).String()]
		idx, err := index(dfImportance, workstyle)
		if err != nil {
			continue
		}
		diff := i - idx
		sumSquared = sumSquared + float64(diff*diff)
	}
	finalScore := 1 - ((6 * sumSquared) / (n * (n*n - 1)))
	match.MatchPercent = (finalScore + 1) * 0.5 * 100
	return match
}

func index(df dataframe.DataFrame, workstyle string) (int, error) {
	for i := 0; i < df.Nrow(); i++ {
		row := df.Subset(i)
		if row.Col("Name").Elem(0).String() == workstyle {
			return i, nil
		}
	}
	return -1, errors.New("did not find matching workstyle")
}

// newAppClient creates new appClient
func newAppClient(app *Application) appClient {
	return appClient{app: app}
}
