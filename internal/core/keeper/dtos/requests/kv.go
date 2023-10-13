//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/models"

	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

// AddKeysRequest defines the Request Content for POST Key DTO.
type AddKeysRequest struct {
	dtoCommon.BaseRequest `json:",inline"`
	Value                 interface{} `json:"value,omitempty"`
}

// Validate checks if the fields are valid of the AddKeysRequest struct
func (a AddKeysRequest) Validate() errors.EdgeX {
	// check if Value field is nil
	if a.Value == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "the value field is undefined", nil)
	}
	// check if Value field is an empty map
	if v, ok := a.Value.(map[string]interface{}); ok {
		if len(v) == 0 {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "the value field is an empty object", nil)
		}
	}

	return nil
}

// AddKeysReqToKVModels transforms the AddKeysRequest DTO to the KV model
func AddKeysReqToKVModels(req AddKeysRequest, key string) models.KV {
	var kv models.KV
	kv.Value = req.Value
	kv.Key = key

	return kv
}
