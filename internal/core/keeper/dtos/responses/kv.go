//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package responses

import (
	kpModels "github.com/edgexfoundry/edgex-go/internal/core/keeper/models"
	dtoCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
)

// KVResponse defines the Response Content for GET KV DTOs.
type KVResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Value                  string `json:"value"`
}

type MultiKVResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Response               []kpModels.KVResponse `json:"response"`
}

func NewMultiKVResponse(requestId string, message string, statusCode int, configs []kpModels.KVResponse) MultiKVResponse {
	return MultiKVResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Response:     configs,
	}
}

// KeysResponse defines the Response Content for the updated or deleted key paths
type KeysResponse struct {
	dtoCommon.BaseResponse `json:",inline"`
	Response               []kpModels.KeyOnly `json:"response"`
}

func NewKeysResponse(requestId string, message string, statusCode int, keys []kpModels.KeyOnly) KeysResponse {
	return KeysResponse{
		BaseResponse: dtoCommon.NewBaseResponse(requestId, message, statusCode),
		Response:     keys,
	}
}
