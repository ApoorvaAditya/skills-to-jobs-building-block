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

package web

import (
	"application/core"
	"application/core/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/go-gota/gota/dataframe"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// ClientAPIsHandler handles the client rest APIs implementation
type ClientAPIsHandler struct {
	app *core.Application
}

func (h ClientAPIsHandler) getExample(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	surveyData, err := h.app.Client.GetExample(claims.OrgID, claims.AppID, id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeExample, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(surveyData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) getOccupationData(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	code := params["code"]
	if len(code) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("code"), nil, http.StatusBadRequest, false)
	}

	occupationData, err := h.app.Client.GetOccupationData(code)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeOccupationData, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(occupationData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) getAllOccupationDatas(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	occupationData, err := h.app.Client.GetAllOccupationDatas()
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeOccupationData, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(occupationData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) getUserMatchingResult(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	id := claims.Subject
	userMatchingResult, err := h.app.Client.GetUserMatchingResult(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeUserMatchingResult, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(userMatchingResult)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) createUserMatchingResult(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	var requestData model.UserMatchingResult
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, logutils.TypeRequestBody, nil, err, http.StatusBadRequest, true)
	}
	requestData.ID = claims.Subject
	userMatchingResult, err := h.app.Client.CreateUserMatchingResult(requestData)
	if err != nil || userMatchingResult == nil {
		return l.HTTPResponseErrorAction(logutils.ActionCreate, model.TypeUserMatchingResult, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(userMatchingResult)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) updateUserMatchingResult(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	var requestData model.UserMatchingResult
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, logutils.TypeRequestBody, nil, err, http.StatusBadRequest, true)
	}

	id := claims.Subject
	requestData.ID = id
	err = h.app.Client.UpdateUserMatchingResult(requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUpdate, model.TypeUserMatchingResult, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

func (h ClientAPIsHandler) deleteUserMatchingResult(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	id := claims.Subject
	err := h.app.Client.DeleteUserMatchingResult(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionDelete, model.TypeUserMatchingResult, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

func (h ClientAPIsHandler) getSurveyData(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	surveyData, err := h.app.Client.GetSurveyData(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionGet, model.TypeSurveyData, nil, err, http.StatusInternalServerError, true)
	}

	response, err := json.Marshal(surveyData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) createSurveyData(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	var requestData model.SurveyData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, logutils.TypeRequestBody, nil, err, http.StatusBadRequest, true)
	}

	surveyData, err := h.app.Client.CreateSurveyData(requestData)
	surveyData.ID = claims.Subject
	if err != nil || surveyData == nil {
		return l.HTTPResponseErrorAction(logutils.ActionCreate, model.TypeSurveyData, nil, err, http.StatusInternalServerError, true)
	}
	go h.runMatchingAlgoAndCreateUserData(*surveyData)

	response, err := json.Marshal(surveyData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionMarshal, logutils.TypeResponseBody, nil, err, http.StatusInternalServerError, false)
	}
	return l.HTTPResponseSuccessJSON(response)
}

func (h ClientAPIsHandler) updateSurveyData(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	var requestData model.SurveyData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUnmarshal, logutils.TypeRequestBody, nil, err, http.StatusBadRequest, true)
	}

	requestData.ID = id
	err = h.app.Client.UpdateSurveyData(requestData)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionUpdate, model.TypeSurveyData, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

func (h ClientAPIsHandler) deleteSurveyData(l *logs.Log, r *http.Request, claims *tokenauth.Claims) logs.HTTPResponse {
	params := mux.Vars(r)
	id := params["id"]
	if len(id) <= 0 {
		return l.HTTPResponseErrorData(logutils.StatusMissing, logutils.TypePathParam, logutils.StringArgs("id"), nil, http.StatusBadRequest, false)
	}

	err := h.app.Client.DeleteSurveyData(id)
	if err != nil {
		return l.HTTPResponseErrorAction(logutils.ActionDelete, model.TypeSurveyData, nil, err, http.StatusInternalServerError, true)
	}

	return l.HTTPResponseSuccess()
}

// NewClientAPIsHandler creates new client API handler instance
func NewClientAPIsHandler(app *core.Application) ClientAPIsHandler {
	return ClientAPIsHandler{app: app}
}

func (h ClientAPIsHandler) runMatchingAlgoAndCreateUserData(surveyData model.SurveyData) {
	occupations, err := h.app.Client.GetAllOccupationDatas()
	if err != nil {
		return
	}

	matches := h.runMatchingAlgo(surveyData.Scores, occupations)
	userMatchingResult := model.UserMatchingResult{
		ID:      surveyData.ID,
		Matches: matches,
		Version: "",
	}

	_, err = h.app.Client.GetUserMatchingResult(userMatchingResult.ID)
	fmt.Println(err)
	if err != nil {
		h.app.Client.CreateUserMatchingResult(userMatchingResult)
	} else {
		h.app.Client.UpdateUserMatchingResult(userMatchingResult)
	}
}

func (h ClientAPIsHandler) runMatchingAlgo(userScores []model.WorkstyleScore, occupations []model.OccupationData) []model.Match {
	matches := make([]model.Match, 0)
	for _, occupation := range occupations {
		workstyles, err := h.app.Client.GetWorkstyleDatasForOccupation(occupation.Code)
		if err != nil {
			return make([]model.Match, 0)
		}
		occupationMatch := h.runMatchingAlgorithmPerOccupation(occupation, userScores, workstyles)
		matches = append(matches, occupationMatch)
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].MatchPercent > matches[j].MatchPercent
	})
	return matches
}

func (h ClientAPIsHandler) runMatchingAlgorithmPerOccupation(occupation model.OccupationData, userScores []model.WorkstyleScore, workstyles []model.WorkstyleData) model.Match {
	match := model.Match{Occupation: occupation}

	df_userScoresUnsorted := dataframe.LoadStructs(userScores)
	df_userScores := df_userScoresUnsorted.Arrange(dataframe.Sort("Score"))
	df_importanceUnsorted := dataframe.LoadStructs(workstyles)
	df_importance := df_importanceUnsorted.Arrange(dataframe.Sort("DataValue"))

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

	sum_squared := 0.0
	n := float64(df_userScores.Nrow())
	for i := 0; i < df_userScores.Nrow(); i++ {
		row := df_userScores.Subset(i)
		workstyle := bessiToWorkstyles[row.Col("Workstyle").Elem(0).String()]
		idx, err := Index(df_importance, workstyle)
		if err != nil {
			return match
		}
		diff := i - idx
		diff_squared := diff * diff
		sum_squared = sum_squared + float64(diff_squared)
	}
	final_score := 1 - ((6 * sum_squared) / (n * (n*n - 1)))
	match.MatchPercent = (final_score + 1) * 0.5 * 100
	return match
}

func Index(df dataframe.DataFrame, workstyle string) (int, error) {
	for i := 0; i < df.Nrow(); i++ {
		row := df.Subset(i)
		if row.Col("ElementName").Elem(0).String() == workstyle {
			return i, nil
		}
	}
	return -1, errors.New("did not find matching workstyle")
}
