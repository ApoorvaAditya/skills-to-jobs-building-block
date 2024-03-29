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
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/v3/authservice"
	"github.com/rokwire/core-auth-library-go/v3/tokenauth"

	"github.com/rokwire/logging-library-go/v2/logs"
	"github.com/rokwire/logging-library-go/v2/logutils"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Adapter entity
type Adapter struct {
	baseURL   string
	port      string
	serviceID string

	auth *Auth

	cachedYamlDoc []byte

	defaultAPIsHandler DefaultAPIsHandler
	clientAPIsHandler  ClientAPIsHandler
	adminAPIsHandler   AdminAPIsHandler

	app *core.Application

	logger *logs.Logger
}

type handlerFunc = func(*logs.Log, *http.Request, *tokenauth.Claims) logs.HTTPResponse

// Start starts the module
func (a Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)

	// handle apis
	baseRouter := router.PathPrefix("/" + a.serviceID).Subrouter()
	baseRouter.PathPrefix("/doc/ui").Handler(a.serveDocUI())
	baseRouter.HandleFunc("/doc", a.serveDoc)
	baseRouter.HandleFunc("/version", a.wrapFunc(a.defaultAPIsHandler.version, nil)).Methods("GET")

	mainRouter := baseRouter.PathPrefix("/api").Subrouter()

	// Client APIs

	// Occupation API
	mainRouter.HandleFunc("/occupation/{code}", a.wrapFunc(a.clientAPIsHandler.getOccupationData, a.auth.client.User)).Methods("GET")
	// mainRouter.HandleFunc("/occupation", a.wrapFunc(a.clientAPIsHandler.getAllOccupationDatas, a.auth.client.User)).Methods("GET")

	// UserMatchingResult API
	mainRouter.HandleFunc("/user-match-results", a.wrapFunc(a.clientAPIsHandler.getUserMatchingResult, a.auth.client.User)).Methods("GET")
	// mainRouter.HandleFunc("/user-match-results", a.wrapFunc(a.clientAPIsHandler.deleteUserMatchingResult, a.auth.client.User)).Methods("DELETE")

	// Survey Data API
	// mainRouter.HandleFunc("/survey-data/{id}", a.wrapFunc(a.clientAPIsHandler.getSurveyData, a.auth.client.User)).Methods("GET")
	mainRouter.HandleFunc("/survey-data", a.wrapFunc(a.clientAPIsHandler.createSurveyData, a.auth.client.User)).Methods("POST")
	// mainRouter.HandleFunc("/survey-data/{id}", a.wrapFunc(a.clientAPIsHandler.updateSurveyData, a.auth.client.User)).Methods("PUT")
	// mainRouter.HandleFunc("/survey-data/{id}", a.wrapFunc(a.clientAPIsHandler.deleteSurveyData, a.auth.client.User)).Methods("DELETE")

	// Admin APIs
	adminRouter := mainRouter.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.adminAPIsHandler.getConfig, a.auth.admin.Permissions)).Methods("GET")
	adminRouter.HandleFunc("/configs", a.wrapFunc(a.adminAPIsHandler.getConfigs, a.auth.admin.Permissions)).Methods("GET")
	adminRouter.HandleFunc("/configs", a.wrapFunc(a.adminAPIsHandler.createConfig, a.auth.admin.Permissions)).Methods("POST")
	adminRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.adminAPIsHandler.updateConfig, a.auth.admin.Permissions)).Methods("PUT")
	adminRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.adminAPIsHandler.deleteConfig, a.auth.admin.Permissions)).Methods("DELETE")

	// BB APIs
	// bbsRouter := mainRouter.PathPrefix("/bbs").Subrouter()

	// TPS APIs
	// tpsRouter := mainRouter.PathPrefix("/tps").Subrouter()

	// System APIs
	// systemRouter := mainRouter.PathPrefix("/system").Subrouter()

	a.logger.Fatalf("Error serving: %v", http.ListenAndServe(":"+a.port, router))
}

func (a Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")

	if a.cachedYamlDoc != nil {
		http.ServeContent(w, r, "", time.Now(), bytes.NewReader([]byte(a.cachedYamlDoc)))
	} else {
		http.ServeFile(w, r, "./driver/web/docs/gen/def.yaml")
	}
}

func (a Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/doc", a.baseURL)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func loadDocsYAML(baseServerURL string) ([]byte, error) {
	data, _ := os.ReadFile("./driver/web/docs/gen/def.yaml")
	yamlMap := yaml.MapSlice{}
	err := yaml.Unmarshal(data, &yamlMap)
	if err != nil {
		return nil, err
	}

	for index, item := range yamlMap {
		if item.Key == "servers" {
			var serverList []interface{}
			if baseServerURL != "" {
				serverList = []interface{}{yaml.MapSlice{yaml.MapItem{Key: "url", Value: baseServerURL}}}
			}

			item.Value = serverList
			yamlMap[index] = item
			break
		}
	}

	yamlDoc, err := yaml.Marshal(&yamlMap)
	if err != nil {
		return nil, err
	}

	return yamlDoc, nil
}

func (a Adapter) wrapFunc(handler handlerFunc, authorization tokenauth.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logObj := a.logger.NewRequestLog(req)

		logObj.RequestReceived()

		var response logs.HTTPResponse
		if authorization != nil {
			responseStatus, claims, err := authorization.Check(req)
			if err != nil {
				logObj.SendHTTPResponse(w, logObj.HTTPResponseErrorAction(logutils.ActionValidate, logutils.TypeRequest, nil, err, responseStatus, true))
				return
			}

			if claims != nil {
				logObj.SetContext("account_id", claims.Subject)
			}
			response = handler(logObj, req, claims)
		} else {
			response = handler(logObj, req, nil)
		}

		logObj.SendHTTPResponse(w, response)
		logObj.RequestComplete()
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(baseURL string, port string, serviceID string, app *core.Application, serviceRegManager *authservice.ServiceRegManager, logger *logs.Logger) Adapter {
	yamlDoc, err := loadDocsYAML(baseURL)
	if err != nil {
		logger.Fatalf("error parsing docs yaml - %s", err.Error())
	}

	auth, err := NewAuth(serviceRegManager)
	if err != nil {
		logger.Fatalf("error creating auth - %s", err.Error())
	}

	defaultAPIsHandler := NewDefaultAPIsHandler(app)
	clientAPIsHandler := NewClientAPIsHandler(app)
	adminAPIsHandler := NewAdminAPIsHandler(app)
	return Adapter{baseURL: baseURL, port: port, serviceID: serviceID, cachedYamlDoc: yamlDoc, auth: auth, defaultAPIsHandler: defaultAPIsHandler,
		clientAPIsHandler: clientAPIsHandler, adminAPIsHandler: adminAPIsHandler, app: app, logger: logger}
}
