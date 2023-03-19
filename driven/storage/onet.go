package storage

import (
	"application/core/model"

	"github.com/rokwire/logging-library-go/v2/errors"
	"github.com/rokwire/logging-library-go/v2/logutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetOnetData finds OnetData by id
func (a Adapter) GetOnetDataById(id string) (*model.OnetData, error) {
	filter := bson.M{"_id": id}

	var data *model.OnetData
	err := a.db.OnetDatas.FindOneWithContext(a.context, filter, &data, nil)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeOnetData, filterArgs(filter), err)
	}

	return data, nil
}

// GetManyOnetData finds any number of OnetData that is specified
func (a Adapter) GetManyOnetData(option *options.FindOptions) (*model.OnetData, error) {
	filter := bson.M{}
	var data *model.OnetData
	err := a.db.OnetDatas.FindWithContext(a.context, filter, &data, option)
	if err != nil {
		return nil, errors.WrapErrorAction(logutils.ActionFind, model.TypeOnetData, filterArgs(filter), err)
	}

	return data, nil
}
