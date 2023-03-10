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
	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"

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
	bbsAPIsHandler     BBsAPIsHandler
	tpsAPIsHandler     TPSAPIsHandler
	systemAPIsHandler  SystemAPIsHandler

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
	mainRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.clientAPIsHandler.getExample, a.auth.client.Permissions)).Methods("GET")

	// Admin APIs
	adminRouter := mainRouter.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.adminAPIsHandler.getExample, a.auth.admin.Permissions)).Methods("GET")
	adminRouter.HandleFunc("/examples", a.wrapFunc(a.adminAPIsHandler.createExample, a.auth.admin.Permissions)).Methods("POST")
	adminRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.adminAPIsHandler.updateExample, a.auth.admin.Permissions)).Methods("PUT")
	adminRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.adminAPIsHandler.deleteExample, a.auth.admin.Permissions)).Methods("DELETE")

	// BB APIs
	bbsRouter := mainRouter.PathPrefix("/bbs").Subrouter()
	bbsRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.bbsAPIsHandler.getExample, a.auth.bbs.Permissions)).Methods("GET")

	// TPS APIs
	tpsRouter := mainRouter.PathPrefix("/tps").Subrouter()
	tpsRouter.HandleFunc("/examples/{id}", a.wrapFunc(a.tpsAPIsHandler.getExample, a.auth.tps.Permissions)).Methods("GET")

	// System APIs
	systemRouter := mainRouter.PathPrefix("/system").Subrouter()
	systemRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.systemAPIsHandler.getConfig, a.auth.system.Permissions)).Methods("GET")
	systemRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.systemAPIsHandler.saveConfig, a.auth.system.Permissions)).Methods("PUT")
	systemRouter.HandleFunc("/configs/{id}", a.wrapFunc(a.systemAPIsHandler.deleteConfig, a.auth.system.Permissions)).Methods("DELETE")

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

			logObj.SetContext("account_id", claims.Subject)
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
	bbsAPIsHandler := NewBBsAPIsHandler(app)
	return Adapter{baseURL: baseURL, port: port, serviceID: serviceID, cachedYamlDoc: yamlDoc, auth: auth, defaultAPIsHandler: defaultAPIsHandler,
		clientAPIsHandler: clientAPIsHandler, adminAPIsHandler: adminAPIsHandler, bbsAPIsHandler: bbsAPIsHandler, app: app, logger: logger}
}
