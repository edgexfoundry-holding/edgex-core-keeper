//
// Copyright (C) 2024 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/config"

	"github.com/edgexfoundry/go-mod-bootstrap/v3/di"
)

// ConfigurationName contains the name of keeper's config.ConfigurationStruct implementation in the DIC.
var ConfigurationName = di.TypeInstanceToName(config.ConfigurationStruct{})

// ConfigurationFrom helper function queries the DIC and returns keeper's config.ConfigurationStruct implementation.
func ConfigurationFrom(get di.Get) *config.ConfigurationStruct {
	return get(ConfigurationName).(*config.ConfigurationStruct)
}
