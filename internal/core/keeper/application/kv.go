//
// Copyright (C) 2024 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/container"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/utils"
	"github.com/edgexfoundry/edgex-go/internal/pkg/correlation"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"
	msgTypes "github.com/edgexfoundry/go-mod-messaging/v3/pkg/types"
)

func Keys(key string, keysOnly bool, isRaw bool, dic *di.Container) (configs []models.KVResponse, err errors.EdgeX) {
	err = utils.ValidateKeys(key)
	if err != nil {
		return configs, errors.NewCommonEdgeXWrapper(err)
	}

	dbClient := container.DBClientFrom(dic.Get)

	configs, err = dbClient.KeeperKeys(key, keysOnly, isRaw)
	if err != nil {
		return configs, errors.NewCommonEdgeXWrapper(err)
	}
	return configs, nil
}

func AddKeys(kv models.KVS, isFlatten bool, dic *di.Container) (keys []models.KeyOnly, err errors.EdgeX) {
	err = utils.ValidateKeys(kv.Key)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	dbClient := container.DBClientFrom(dic.Get)
	keys, err = dbClient.AddKeeperKeys(kv, isFlatten)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}
	return keys, nil
}

func DeleteKeys(key string, prefixMatch bool, dic *di.Container) (keys []models.KeyOnly, err errors.EdgeX) {
	err = utils.ValidateKeys(key)
	if err != nil {
		return keys, errors.NewCommonEdgeXWrapper(err)
	}

	dbClient := container.DBClientFrom(dic.Get)
	keys, err = dbClient.DeleteKeeperKeys(key, prefixMatch)
	if err != nil {
		return keys, errors.NewCommonEdgeXWrapper(err)
	}
	return keys, nil
}

// PublishKeyChange publishes any key value changes in the format of []byte through MessageClient
func PublishKeyChange(data []byte, key string, ctx context.Context, dic *di.Container) {
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)
	msgClient := bootstrapContainer.MessagingClientFrom(dic.Get)
	configuration := container.ConfigurationFrom(dic.Get)
	correlationId := correlation.FromContext(ctx)

	publishTopic := configuration.MessageBus.BaseTopicPrefix + "/" + key
	lc.Debugf("Publishing keeper key change to message queue. Topic: %s; %s: %s", publishTopic, common.CorrelationHeader, correlationId)

	msgEnvelope := msgTypes.NewMessageEnvelope(data, ctx)

	// ensure the message envelope content-type is application/json
	msgEnvelope.ContentType = common.ContentTypeJSON

	err := msgClient.Publish(msgEnvelope, publishTopic)
	if err != nil {
		lc.Errorf("Unable to send message for Key: %s. Correlation-id: %s, Error: %v", key, correlationId, err)
	} else {
		lc.Debugf("Keeper key change published on message queue. Topic: %s, Correlation-id: %s", publishTopic, correlationId)
	}
}
