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

package storage

import (
	"application/core/interfaces"
	"context"
	"time"

	"github.com/rokwire/logging-library-go/v2/logs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client
	logger   *logs.Logger

	configs         *collectionWrapper
	occupationData  *collectionWrapper
	matchResults    *collectionWrapper
	surveyResponses *collectionWrapper

	listeners []interfaces.StorageListener
}

func (d *database) start() error {

	d.logger.Info("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(d.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), d.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), d.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(d.mongoDBName)

	configs := &collectionWrapper{database: d, coll: db.Collection("configs")}
	err = d.applyConfigsChecks(configs)
	if err != nil {
		return err
	}

	occupationData := &collectionWrapper{database: d, coll: db.Collection("occupation_data")}
	err = d.applyOccupationDataChecks(occupationData)
	if err != nil {
		return err
	}

	matchResults := &collectionWrapper{database: d, coll: db.Collection("match_results")}
	err = d.applyMatchResultsChecks(matchResults)
	if err != nil {
		return err
	}

	surveyResponses := &collectionWrapper{database: d, coll: db.Collection("survey_responses")}
	err = d.applySurveyResponsesChecks(surveyResponses)
	if err != nil {
		return err
	}

	//assign the db, db client and the collections
	d.db = db
	d.dbClient = client

	d.configs = configs
	d.occupationData = occupationData
	d.matchResults = matchResults
	d.surveyResponses = surveyResponses

	go d.configs.Watch(nil, d.logger)

	return nil
}

func (d *database) applyConfigsChecks(configs *collectionWrapper) error {
	d.logger.Info("apply configs checks.....")

	err := configs.AddIndex(nil, bson.D{primitive.E{Key: "type", Value: 1}, primitive.E{Key: "app_id", Value: 1}, primitive.E{Key: "org_id", Value: 1}}, true)
	if err != nil {
		return err
	}

	d.logger.Info("apply configs passed")
	return nil
}

func (d *database) applyOccupationDataChecks(occupationData *collectionWrapper) error {
	d.logger.Info("apply occupationData checks.....")

	err := occupationData.AddIndex(nil, bson.D{primitive.E{Key: "code", Value: 1}}, true)
	if err != nil {
		return err
	}

	d.logger.Info("apply occupationData passed")
	return nil
}

func (d *database) applyMatchResultsChecks(matchResults *collectionWrapper) error {
	d.logger.Info("apply matchResults checks.....")

	d.logger.Info("apply matchResults passed")
	return nil
}

func (d *database) applySurveyResponsesChecks(surveyResponses *collectionWrapper) error {
	d.logger.Info("apply surveyResponses checks.....")

	d.logger.Info("apply surveyResponses passed")
	return nil
}

func (d *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	d.logger.Infof("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	switch coll {
	case "configs":
		d.logger.Info("configs collection changed")

		for _, listener := range d.listeners {
			go listener.OnConfigsUpdated()
		}
	}
}
