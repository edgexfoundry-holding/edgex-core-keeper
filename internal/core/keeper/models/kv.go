//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package models

import commonDTOs "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"

type KVResponse interface {
	SetKey(string)
}

type KV struct {
	Key string `json:"key,omitempty"`
	StoredData
}

type StoredData struct {
	commonDTOs.DBTimestamp
	Value interface{} `json:"value,omitempty"`
}

type KeyOnly string

func (kv *KV) SetKey(newKey string) {
	kv.Key = newKey
}

func (key *KeyOnly) SetKey(newKey string) {
	*key = KeyOnly(newKey)
}
