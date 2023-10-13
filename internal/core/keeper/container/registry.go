//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/infrastructure/interfaces"
)

// RegistryInterfaceName contains the name of the interfaces.Registry implementation in the DIC.
var RegistryInterfaceName = di.TypeInstanceToName((*interfaces.Registry)(nil))

// RegistryFrom helper function queries the DIC and returns the interfaces.Registry implementation.
func RegistryFrom(get di.Get) interfaces.Registry {
	return get(RegistryInterfaceName).(interfaces.Registry)
}
