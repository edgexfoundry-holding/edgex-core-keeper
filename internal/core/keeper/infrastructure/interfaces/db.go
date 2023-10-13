//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package interfaces

import (
	kpModels "github.com/edgexfoundry/edgex-go/internal/core/keeper/models"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

type DBClient interface {
	KeeperKeys(key string, keyOnly bool, isRaw bool) ([]kpModels.KVResponse, errors.EdgeX)
	AddKeeperKeys(kv kpModels.KV, isFlatten bool) ([]kpModels.KeyOnly, errors.EdgeX)
	DeleteKeeperKeys(key string, isRecurse bool) ([]kpModels.KeyOnly, errors.EdgeX)

	AddRegistration(r kpModels.Registration) (kpModels.Registration, errors.EdgeX)
	DeleteRegistrationByServiceId(id string) errors.EdgeX
	Registrations() ([]kpModels.Registration, errors.EdgeX)
	RegistrationByServiceId(id string) (kpModels.Registration, errors.EdgeX)
	UpdateRegistration(r kpModels.Registration) errors.EdgeX
}
