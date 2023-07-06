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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"
)

// ClientAPIsHandler handles the client rest APIs implementation
type ClientAPIsHandler struct {
	app *core.Application
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
	if err != nil || surveyData == nil {
		return l.HTTPResponseErrorAction(logutils.ActionCreate, model.TypeSurveyData, nil, err, http.StatusInternalServerError, true)
	}
	go h.app.Client.MatchOccupations(*surveyData, claims.Subject)

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
